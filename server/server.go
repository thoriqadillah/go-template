package server

import (
	"app/env"
	"app/lib/log"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/riverqueue/river"
	"github.com/stephenafamo/bob"
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
	Db       *bob.DB
	Queue    *river.Client[pgx.Tx]
}

type Factory func(app *App) Service

var factories = make([]Factory, 0)

func Register(f ...Factory) {
	factories = append(factories, f...)
}

func Create(echo *echo.Echo, db *bob.DB, river *river.Client[pgx.Tx]) *App {
	app := &App{
		echo:     echo,
		Db:       db,
		Queue:    river,
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

	hardctx, hardcancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer hardcancel()

	go func() {
		err := a.Queue.StopAndCancel(hardctx)
		if err != nil && errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
			panic(err)
		}

		if err == nil || errors.Is(err, context.Canceled) {
			logger.Info("Queue stoped")
		}
	}()

	if err := a.echo.Shutdown(hardctx); err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info("Server shut down")
}
