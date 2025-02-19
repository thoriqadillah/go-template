package storage

import (
	"app/lib/storage"
	"app/server"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type storageService struct {
	storage storage.Storage
}

func CreateService(app *server.App) server.Service {
	return &storageService{
		storage: storage.New("local"),
	}
}

func (s *storageService) serve(c echo.Context) error {
	filename := c.Param("filename")
	ext := filepath.Ext(filename)
	mimetype := mime.TypeByExtension(ext)

	file, err := s.storage.Serve(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return c.Stream(http.StatusOK, mimetype, file)
}

func (s *storageService) CreateRoutes(app *echo.Echo) {
	router := app.Group("storage")

	router.GET("/:filename", s.serve)
}
