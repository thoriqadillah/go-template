package env

import (
	"app/common"
	"os"
)

func env(key string) common.Parser {
	return common.Parse(os.Getenv(key))
}

var (
	AppName       = env("APP_NAME").String("Echo")
	Env           = env("ENV").String("development")
	Dev           = env("ENV").String("development") == "development"
	Prod          = env("ENV").String("production") == "production"
	CorsOrigin    = env("CORS_ORIGIN").String("(.*?)")
	Port          = env("PORT").String(":8080")
	EmailSender   = env("EMAIL_SENDER").String("app@example.com")
	EmailHost     = env("EMAIL_HOST").String("smtp.example.com")
	EmailPort     = env("EMAIL_PORT").Int(587)
	EmailSecure   = env("EMAIL_SECURE").Bool(false)
	EmailUsername = env("EMAIL_USERNAME").String("app@example.com")
	EmailPassword = env("EMAIL_PASSWORD").String("password")
	StorageDriver = env("STORAGE_DRIVER").String("local")
)
