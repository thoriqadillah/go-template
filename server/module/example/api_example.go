package example

import (
	"app/common"
	"app/env"
	"app/lib/auth"
	"app/lib/notifier"
	"app/lib/storage"
	"app/server"
	"net/http"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type exampleService struct {
	storage storage.Storage
	email   notifier.Notifier
}

func CreateService(app *server.App) server.Service {
	return &exampleService{
		storage: storage.New(env.STORAGE_DRIVER),
		email:   notifier.New(notifier.EmailNotifier, notifier.WithRiverQueue(app.Queue)),
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
	err := s.email.Send(notifier.Message{
		Subject: "Test Email",
		Body:    "example.html",
		Data: notifier.Data{
			"message": "Hello World",
		},
	})

	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (s *exampleService) example(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func (s *exampleService) paginate(c echo.Context) error {
	params := c.QueryParams()
	paginator := common.Paginate(params)

	return c.JSON(http.StatusOK, paginator.CreatePager(nil, 0))
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

	id, _ := uuid.NewV7()
	ext := filepath.Ext(file.Filename)
	filename := id.String() + ext
	url, err := s.storage.Upload(filename, src)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, url)
}

func (s *exampleService) restricted(c echo.Context) error {
	user := auth.User(c)
	return c.JSON(http.StatusOK, echo.Map{
		"id": user.UserId,
	})
}

func (s *exampleService) CreateRoutes(e *echo.Echo) {
	r := e.Group("/api/example")

	r.GET("/", s.example)
	r.GET("/email", s.sendEmail)
	r.POST("/upload", s.uploadFile)
	r.POST("/validate", s.validate)
	r.GET("/paginate", s.paginate)

	restricted := r.Group("/restricted")

	restricted.Use(auth.AuthenticatedMw)
	restricted.GET("/", s.restricted)
}
