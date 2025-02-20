package account

import (
	"app/lib/auth"
	"app/server"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type accountService struct{}

func CreateService(app *server.App) server.Service {
	return &accountService{}
}

func (s *accountService) login(c echo.Context) error {
	// TODO: perform login

	id := uuid.NewString()
	token, err := auth.SignToken(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}

func (s *accountService) CreateRoutes(echo *echo.Echo) {
	r := echo.Group("/account")

	r.POST("/login", s.login)
}
