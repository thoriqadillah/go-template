package example

import (
	"app/lib/notification/notifier"
	"net/http"

	"github.com/labstack/echo/v4"
)

type exampleService struct{}

func CreateService() *exampleService {
	return &exampleService{}
}

func (s *exampleService) Init() {
	// Do something when the service is initialized
}

func (s *exampleService) Close() {
	// Do something when the service is closed
}

func (s *exampleService) sendEmail(c echo.Context) error {
	notifier.New("email").Send(notifier.Send{
		Subject: "Test Email",
		Message: "This is a test email",
	})

	return c.String(http.StatusOK, "Email sent")
}

func (s *exampleService) example(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func (s *exampleService) CreateRoutes(app *echo.Echo) {
	router := app.Group("example")

	router.GET("/", s.example)
	router.GET("/email", s.sendEmail)
}
