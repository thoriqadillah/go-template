package example

import (
	"app/env"
	"app/lib/notification/notifier"
	"app/lib/storage"
	"net/http"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type exampleService struct {
	storage storage.Storage
}

func CreateService() *exampleService {
	return &exampleService{
		storage: storage.New(env.StorageDriver),
	}
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
	err := notifier.New("email").Send(notifier.Send{
		Subject:  "Test Email",
		Template: "example.html",
	})

	if err != nil {
		return err
	}

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

	ext := filepath.Ext(file.Filename)
	filename := uuid.NewString() + ext
	url, err := s.storage.Upload(filename, src)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, url)
}

func (s *exampleService) CreateRoutes(app *echo.Echo) {
	router := app.Group("example")

	router.GET("/", s.example)
	router.GET("/email", s.sendEmail)
	router.POST("/upload", s.uploadFile)
	router.POST("/validate", s.validate)
}
