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

func (s *accountService) getUser(c echo.Context) error {
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

func (s *accountService) CreateRoutes(echo *echo.Echo) {
	r := echo.Group("/account")

	r.POST("/login", s.login)
	r.POST("/signup", s.signup)
	r.GET("/", s.getUser, auth.Middleware())
}
