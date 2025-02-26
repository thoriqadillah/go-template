package main

import (
	"app/env"
	"app/lib/log"
	"app/lib/validator"
	"app/server"
	_ "app/server/module"
	"context"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"

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

	echo := echo.New()
	echo.Binder = &customBinder{}

	logger := log.Logger()
	defer logger.Sync()

	echo.Use(middleware.Recover())
	echo.Use(log.Middleware())
	echo.Use(middleware.Gzip())
	echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			re := regexp.MustCompile(env.CORS_ORIGIN)
			return re.MatchString(origin), nil
		},
		AllowCredentials: true,
	}))

	server.Run(ctx, echo)
}
