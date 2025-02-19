package server

import (
	"app/env"
	"app/lib/log"
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

var logger = log.Logger()

type Service interface {
	CreateRoutes(app *echo.Echo)
}

type Initter interface {
	Init()
}

type Closer interface {
	Close()
}

type App struct {
	echo     *echo.Echo
	services []Service
}

type Factory func(app *App) Service

var factories = make([]Factory, 0)

func Register(f ...Factory) {
	factories = append(factories, f...)
}

func Create(echo *echo.Echo) *App {
	app := &App{
		echo:     echo,
		services: make([]Service, 0),
	}

	for _, factory := range factories {
		service := factory(app)
		service.CreateRoutes(echo)

		if initter, ok := service.(Initter); ok {
			initter.Init()
		}

		app.services = append(app.services, service)
	}

	return app
}

func (a *App) Start(ctx context.Context) {
	go func() {
		err := a.echo.Start(env.PORT)
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal(err.Error())
		}
	}()

	<-ctx.Done()
	logger.Info("Interrupt signal received. Shutting down")
	for _, service := range a.services {
		if closer, ok := service.(Closer); ok {
			closer.Close()
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.echo.Shutdown(ctx); err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info("Server shut down")
}
