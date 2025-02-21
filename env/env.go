package env

import (
	"app/common"
	"os"
)

func env(key string) common.Parser {
	return common.Parse(os.Getenv(key))
}

var (
	APP_NAME       = env("APP_NAME").String("App")
	ENV            = env("ENV").String("development")
	DEV            = env("ENV").String("development") == "development"
	PROD           = env("ENV").String("production") == "production"
	CORS_ORIGIN    = env("CORS_ORIGIN").String("(.*?)")
	JWT_SECRET     = env("JWT_SECRET").String("secret")
	PORT           = env("PORT").String(":8080")
	DB_URL         = env("DB_URL").String("postgresql://postgres@localhost:5432/app?sslmode=disable")
	REDIS_URL      = env("REDIS_URL").String("redis://:@localhost:6379")
	EMAIL_SENDER   = env("EMAIL_SENDER").String("hallo@app.com")
	EMAIL_HOST     = env("EMAIL_HOST").String("smtp.app.com")
	EMAIL_PORT     = env("EMAIL_PORT").Int(587)
	EMAIL_SECURE   = env("EMAIL_SECURE").Bool(false)
	EMAIL_USERNAME = env("EMAIL_USERNAME").String("hallo@app.com")
	EMAIL_PASSWORD = env("EMAIL_PASSWORD").String("password")
	STORAGE_DRIVER = env("STORAGE_DRIVER").String("local")
)
