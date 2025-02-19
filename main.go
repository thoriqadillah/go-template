package main

import (
	"app/db"
	"app/env"
	"app/lib/log"
	"app/lib/validator"
	"app/queue"
	"app/server"
	_ "app/server/module"
	"context"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/brpaz/echozap"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type customBinder struct {
	echo.DefaultBinder
}

func (b *customBinder) Bind(i interface{}, c echo.Context) error {
	if err := b.DefaultBinder.Bind(i, c); err != nil {
		return err
	}

	validate, ok := i.(validator.Validator)
	if !ok {
		return nil
	}

	if err := validate.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, validator.Translate(err))
	}

	return nil
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGABRT, syscall.SIGTERM)
	defer stop()

	godotenv.Load()

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

	echo := echo.New()
	echo.Binder = &customBinder{}

	logger := log.Logger()
	defer logger.Sync()

	echo.Use(middleware.Recover())
	echo.Use(echozap.ZapLogger(logger))
	echo.Use(middleware.Gzip())
	echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			re := regexp.MustCompile(env.CORS_ORIGIN)
			return re.MatchString(origin), nil
		},
	}))

	app := server.Create(echo, db, river)
	app.Start(ctx)
}
