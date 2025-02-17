package main

import (
	"app/env"
	"app/lib/log"
	"app/lib/validator"
	"app/server"
	"net/http"
	"regexp"

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
	godotenv.Load()
	echo := echo.New()
	echo.Binder = &customBinder{}

	logger := log.Logger()

	echo.Use(middleware.Recover())
	echo.Use(echozap.ZapLogger(logger))
	echo.Use(middleware.Gzip())
	echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			re := regexp.MustCompile(env.CorsOrigin)
			return re.MatchString(origin), nil
		},
	}))

	app := server.Create(echo)
	app.Start()
}
