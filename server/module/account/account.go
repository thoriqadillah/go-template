package account

import (
	"app/lib/auth"
	"app/lib/notifier"
	"app/server"
	"context"
	"crypto/rand"
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
		store: NewStore(app.BobDB),
		cache: app.Redis,
		emailer: notifier.New(notifier.EmailNotifier,
			notifier.WithRiverQueue(app.RiverQueue),
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

	if err := s.cache.Set(ctx, email, otp, time.Minute*30).Err(); err != nil {
		return err
	}

	return s.emailer.Send(notifier.Message{
		To:       []string{email},
		Subject:  "Email Verification",
		Template: "verify.html",
		Data: notifier.Data{
			"otp": otp,
		},
	})
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

	token, err := auth.SignToken(user.ID.String())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"token": token,
	})
}

func (s *accountService) loginOauth(c echo.Context) error {
	// TODO
	return c.NoContent(http.StatusOK)
}

func (s *accountService) signupOauth(c echo.Context) error {
	// TODO
	return c.NoContent(http.StatusOK)
}

func (s *accountService) user(c echo.Context) error {
	ctx := c.Request().Context()

	claims := auth.User(c)
	user, err := s.store.Get(ctx, claims.UserId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"id":         user.ID,
		"email":      user.Email,
		"name":       user.Name,
		"source":     user.Source,
		"createdAt":  user.CreatedAt,
		"updatedAt":  user.UpdatedAt,
		"verifiedAt": user.VerifiedAt,
	})
}

func (s *accountService) auth(c echo.Context) error {
	claims := auth.User(c)
	return c.JSON(http.StatusOK, echo.Map{
		"id": claims.UserId,
	})
}

func (s *accountService) logout(c echo.Context) error {
	// TODO
	return c.NoContent(http.StatusOK)
}

func (s *accountService) refreshToken(c echo.Context) error {
	// TODO
	return c.NoContent(http.StatusOK)
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

	otp, err := s.cache.Get(ctx, user.Email).Result()
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

	if err := s.cache.Del(ctx, user.Email).Err(); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (s *accountService) resetPassword(c echo.Context) error {
	// TODO
	return c.NoContent(http.StatusOK)
}

func (s *accountService) changePassword(c echo.Context) error {
	// TODO
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
}
