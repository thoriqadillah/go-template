package main

import (
	"app/env"
	"app/server"
	"regexp"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	godotenv.Load()
	echo := echo.New()

	echo.Use(middleware.Recover())
	echo.Use(middleware.Logger())
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
