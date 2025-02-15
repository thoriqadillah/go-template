package example

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type exampleService struct{}

func CreateService() *exampleService {
	return &exampleService{}
}

func (s *exampleService) example(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func (s *exampleService) CreateRoutes(app *echo.Echo) {
	app.GET("/", s.example)
}
