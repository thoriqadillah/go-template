package account

import (
	"app/lib/auth"
	"app/lib/notifier"
	"app/server"
	"net/http"

	"github.com/labstack/echo/v4"
)

type accountService struct {
	store   Store
	emailer notifier.Notifier
}

func CreateService(app *server.App) server.Service {
	return &accountService{
		store:   NewStore(app.Db),
		emailer: notifier.New(notifier.EmailNotifier, notifier.WithQueue(app.Queue)),
	}
}

func (s *accountService) login(c echo.Context) error {
	ctx := c.Request().Context()

	var user loginUser
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

	// TODO: generate OTP
	err = s.emailer.Send(notifier.Message{
		To:       []string{user.Email},
		Subject:  "Email Verification",
		Template: "verify.html",
		Data:     notifier.Data{
			// TODO
		},
	})

	if err != nil {
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

func (s *accountService) resendVerification(c echo.Context) error {
	ctx := c.Request().Context()

	claims := auth.User(c)
	user, err := s.store.Get(ctx, claims.UserId)
	if err != nil {
		return err
	}

	// TODO: generate OTP
	err = s.emailer.Send(notifier.Message{
		To:       []string{user.Email},
		Subject:  "Email Verification",
		Template: "verify.html",
		Data:     notifier.Data{
			// TODO
		},
	})

	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (s *accountService) verifyUser(c echo.Context) error {
	// TODO
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

func (s *accountService) CreateRoutes(echo *echo.Echo) {
	a := echo.Group("/auth")
	a.POST("/login", s.login)
	a.POST("/signup", s.signup)

	acc := echo.Group("/account", auth.Middleware())
	acc.GET("/", s.auth)
	acc.GET("/user", s.user)
	acc.POST("/logout", s.logout)
	acc.POST("/refresh-token", s.refreshToken)
	acc.POST("/reset-password", s.resetPassword)
	acc.POST("/change-password", s.changePassword)
	acc.POST("/verify", s.verifyUser)
	acc.POST("/resend-verification", s.resendVerification)
}
