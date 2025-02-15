package main

import (
	"app/env"
	"fmt"
	"regexp"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	godotenv.Load()
	app := echo.New()

	app.Use(middleware.Recover())
	app.Use(middleware.Logger())
	app.Use(middleware.Gzip())
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			re := regexp.MustCompile(env.CorsOrigin)
			return re.MatchString(origin), nil
		},
	}))

	app.Start(fmt.Sprintf(":%d", env.Port))
}
