package account

import (
	"app/env"
	"app/lib/auth"
	"app/lib/auth/oauth"
	"app/lib/notifier"
	"app/server"
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type accountService struct {
	store   Store
	cache   *redis.Client
	emailer notifier.Notifier
}

func CreateService(app *server.App) server.Service {
	return &accountService{
		store: NewStore(app.Db),
		cache: app.Cache,
		emailer: notifier.New(notifier.EmailNotifier,
			notifier.WithRiverQueue(app.Queue),
		),
	}
}

func (s *accountService) generateOTP() (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()) + 100000, nil
}

func (s *accountService) sendVerificationEmail(ctx context.Context, email string) error {
	otp, err := s.generateOTP()
	if err != nil {
		return err
	}

	key := fmt.Sprintf("verify:%s", email)
	if err := s.cache.Set(ctx, key, otp, time.Minute*30).Err(); err != nil {
		return err
	}

	return s.emailer.Send(notifier.Message{
		To:      []string{email},
		Subject: "Email Verification",
		Body:    "verify.html",
		Data: notifier.Data{
			"otp": otp,
		},
	})
}

func (s *accountService) createRefreshToken(c echo.Context, token string) {
	sameSite := http.SameSiteNoneMode
	if c.Scheme() == "https" {
		sameSite = http.SameSiteLaxMode
	}

	cookie := http.Cookie{
		Name:        "refresh-token",
		HttpOnly:    true,
		SameSite:    sameSite,
		Secure:      c.Scheme() == "https",
		Partitioned: c.Scheme() == "https",
		Path:        "/",
		Expires:     time.Now().Add(7 * 24 * time.Hour),
		Value:       token,
	}

	c.SetCookie(&cookie)
}

func (s *accountService) login(c echo.Context) error {
	ctx := c.Request().Context()

	var user login
	if err := c.Bind(&user); err != nil {
		return err
	}

	id, err := s.store.Login(ctx, user.Email, user.Password)
	if err != nil {
		return err
	}

	token, err := auth.SignToken(id)
	if err != nil {
		return err
	}

	s.createRefreshToken(c, token)
	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}

func (s *accountService) signup(c echo.Context) error {
	ctx := c.Request().Context()

	var data createUser
	if err := c.Bind(&data); err != nil {
		return err
	}

	user, err := s.store.Signup(ctx, data)
	if err != nil {
		return err
	}

	if err = s.sendVerificationEmail(ctx, user.Email); err != nil {
		return err
	}

	token, err := auth.SignToken(user.Id)
	if err != nil {
		return err
	}

	s.createRefreshToken(c, token)
	return c.JSON(http.StatusCreated, echo.Map{
		"token": token,
	})
}

func (s *accountService) loginOauth(c echo.Context) error {
	ctx := c.Request().Context()

	var req oauthReq
	if err := c.Bind(&req); err != nil {
		return err
	}

	oauth := oauth.Create(req.Provider)
	user, err := oauth.Validate(ctx, req.Token)
	if err != nil {
		return err
	}

	id, err := s.store.Login(ctx, user.Email)
	if err != nil {
		return err
	}

	token, err := auth.SignToken(id)
	if err != nil {
		return err
	}

	s.createRefreshToken(c, token)
	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}

func (s *accountService) signupOauth(c echo.Context) error {
	ctx := c.Request().Context()

	var req oauthReq
	if err := c.Bind(&req); err != nil {
		return err
	}

	oauth := oauth.Create(req.Provider)
	user, err := oauth.Validate(ctx, req.Token)
	if err != nil {
		return err
	}

	data := createUser{
		Name:   user.Name,
		Email:  user.Email,
		Source: req.Provider,
	}

	u, err := s.store.Signup(ctx, data)
	if err != nil {
		return err
	}

	token, err := auth.SignToken(u.Id)
	if err != nil {
		return err
	}

	s.createRefreshToken(c, token)
	return c.JSON(http.StatusCreated, echo.Map{
		"token": token,
	})
}

func (s *accountService) user(c echo.Context) error {
	ctx := c.Request().Context()

	claims := auth.User(c)
	user, err := s.store.Get(ctx, claims.UserId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (s *accountService) auth(c echo.Context) error {
	claims := auth.User(c)
	return c.JSON(http.StatusOK, echo.Map{
		"id": claims.UserId,
	})
}

func (s *accountService) deleteCookie(c echo.Context) {
	sameSite := http.SameSiteNoneMode
	if env.PROD {
		sameSite = http.SameSiteLaxMode
	}

	cookie := http.Cookie{
		Name:        "refresh-token",
		HttpOnly:    true,
		SameSite:    sameSite,
		Secure:      env.PROD,
		Partitioned: env.PROD,
		Path:        "/",
		MaxAge:      -1,
		Value:       "",
	}

	c.SetCookie(&cookie)
}

func (s *accountService) logout(c echo.Context) error {
	s.deleteCookie(c)
	return c.NoContent(http.StatusOK)
}

func (s *accountService) refreshToken(c echo.Context) error {
	cookie, err := c.Cookie("refresh-token")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	if cookie == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	claims, err := auth.DecodeToken(cookie.Value)
	if err != nil {
		return err
	}

	token, err := auth.SignToken(claims.UserId)
	if err != nil {
		return err
	}

	s.createRefreshToken(c, token)
	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}

func (s *accountService) sendVerification(c echo.Context) error {
	ctx := c.Request().Context()

	claims := auth.User(c)
	user, err := s.store.Get(ctx, claims.UserId)
	if err != nil {
		return err
	}

	if err = s.sendVerificationEmail(ctx, user.Email); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (s *accountService) verifyUser(c echo.Context) error {
	ctx := c.Request().Context()

	otpquery := c.QueryParam("otp")

	claims := auth.User(c)
	user, err := s.store.Get(ctx, claims.UserId)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("verify:%s", user.Email)
	otp, err := s.cache.Get(ctx, key).Result()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid or expired OTP")
	}

	if otpquery != otp {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid or expired OTP")
	}

	data := updateUser{VerifiedAt: time.Now()}
	if err := s.store.Update(ctx, claims.UserId, data); err != nil {
		return err
	}

	if err := s.cache.Del(ctx, key).Err(); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (s *accountService) resetPassword(c echo.Context) error {
	ctx := c.Request().Context()

	email := c.QueryParam("email")
	if email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Email query is required")
	}

	otp, err := s.generateOTP()
	if err != nil {
		return err
	}

	key := fmt.Sprintf("reset:%s", email)
	if err := s.cache.Set(ctx, key, otp, time.Minute*30).Err(); err != nil {
		return err
	}

	err = s.emailer.Send(notifier.Message{
		Subject: "Password Reset",
		Body:    "reset-password.html",
		To:      []string{email},
		Data: notifier.Data{
			"otp": otp,
		},
	})

	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (s *accountService) changePassword(c echo.Context) error {
	ctx := c.Request().Context()

	var payload changePassword
	if err := c.Bind(&payload); err != nil {
		return err
	}

	key := fmt.Sprintf("reset:%s", payload.Email)
	otp, err := s.cache.Get(ctx, key).Result()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid or expired OTP")
	}

	if payload.Otp != otp {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid or expired OTP")
	}

	user, err := s.store.GetByEmail(ctx, payload.Email)
	if err != nil {
		return err
	}

	update := updateUser{
		Password: payload.Password,
	}
	if err := s.store.Update(ctx, user.Id, update); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (s *accountService) deleteAccount(c echo.Context) error {
	ctx := c.Request().Context()

	user := auth.User(c)
	if err := s.store.Delete(ctx, user.UserId); err != nil {
		return err
	}

	s.deleteCookie(c)
	return c.NoContent(http.StatusOK)
}

func (s *accountService) CreateRoutes(e *echo.Echo) {
	r := e.Group("/api")

	oa := r.Group("/oauth")
	oa.POST("/login", s.loginOauth)
	oa.POST("/signup", s.signupOauth)

	au := r.Group("/auth")
	au.POST("/login", s.login)
	au.POST("/signup", s.signup)
	au.POST("/reset-password", s.resetPassword)
	au.GET("/refresh-token", s.refreshToken)
	au.PATCH("/change-password", s.changePassword)

	acc := r.Group("/account", auth.AuthenticatedMw)
	acc.GET("/", s.auth)
	acc.GET("/user", s.user)
	acc.POST("/logout", s.logout)
	acc.PATCH("/verify", s.verifyUser)
	acc.POST("/send-verification", s.sendVerification)
	acc.DELETE("/", s.deleteAccount)
}
