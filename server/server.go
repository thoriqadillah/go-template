package server

import (
	"app/db"
	"app/env"
	"app/lib/log"
	"app/lib/queue"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/riverqueue/river"
	"github.com/uptrace/bun"
)

var logger = log.Logger()

type Service interface {
	CreateRoutes(e *echo.Echo)
}

type Initter interface {
	Init()
}

type Closer interface {
	Close()
}

type App struct {
	Db    *bun.DB
	Queue *river.Client[pgx.Tx]
	Cache *redis.Client
}

type Factory func(app *App) Service

var factories = make([]Factory, 0)

func Register(f ...Factory) {
	factories = append(factories, f...)
}

func Run(ctx context.Context, echo *echo.Echo) {
	ropt, err := redis.ParseURL(env.REDIS_URL)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(ropt)
	defer rdb.Close()

	db, pool, err := db.Connect(ctx, env.DB_URL)
	if err != nil {
		panic(err)
	}

	defer db.Close()
	defer pool.Close()

	river, err := queue.Start(ctx, pool)
	if err != nil {
		panic(err)
	}

	if err := river.Start(ctx); err != nil {
		panic(err)
	}

	services := make([]Service, 0)

	app := &App{
		Db:    db,
		Queue: river,
		Cache: rdb,
	}

	for _, factory := range factories {
		service := factory(app)
		service.CreateRoutes(echo)

		if initter, ok := service.(Initter); ok {
			initter.Init()
		}

		services = append(services, service)
	}

	go func() {
		err := echo.Start(env.PORT)
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal(err.Error())
		}
	}()

	<-ctx.Done()
	logger.Info("Interrupt signal received. Shutting down")
	for _, service := range services {
		if closer, ok := service.(Closer); ok {
			closer.Close()
		}
	}

	hardctx, hardcancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer hardcancel()

	go func() {
		err := river.StopAndCancel(hardctx)
		if err != nil && errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
			logger.Fatal(err.Error())
		}

		if err == nil || errors.Is(err, context.Canceled) {
			logger.Info("Queue stoped")
		}
	}()

	if err := echo.Shutdown(hardctx); err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info("Server shut down")
}
