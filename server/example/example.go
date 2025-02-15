package example

import (
	"app/lib/notification/notifier"
	"app/lib/storage"
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

func (s *exampleService) validate(c echo.Context) error {
	// try to pass a json with and without required field
	var foo Foo
	if err := c.Bind(&foo); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, foo)
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

func (s *exampleService) uploadFile(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}

	err = storage.New("local").Upload(file.Filename, src)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, "File uploaded")
}

func (s *exampleService) CreateRoutes(app *echo.Echo) {
	router := app.Group("example")

	router.GET("/", s.example)
	router.GET("/email", s.sendEmail)
	router.POST("/upload", s.uploadFile)
	router.POST("/validate", s.validate)
}
