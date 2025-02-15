package server

import (
	"app/env"
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
)

type Service interface {
	CreateRoutes(app *echo.Echo)
}

type Initter interface {
	Init()
}

type Closer interface {
	Close()
}

var services = make([]Service, 0)

func Register(s ...Service) {
	services = append(services, s...)
}

type app struct {
	echo *echo.Echo
}

func Create(echo *echo.Echo) *app {
	for _, service := range services {
		service.CreateRoutes(echo)

		if initter, ok := service.(Initter); ok {
			initter.Init()
		}
	}

	return &app{echo}
}

func (a *app) Start() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		err := a.echo.Start(env.Port)
		if err != nil && err != http.ErrServerClosed {
			a.echo.Logger.Fatal("Shutting down the server...")
		}
	}()

	<-ctx.Done()
	for _, service := range services {
		if closer, ok := service.(Closer); ok {
			closer.Close()
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.echo.Shutdown(ctx); err != nil {
		a.echo.Logger.Fatal(err)
	}
}
