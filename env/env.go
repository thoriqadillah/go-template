package env

import (
	"app/common"
	"os"
)

func env(key string) common.Parser {
	return common.Parse(os.Getenv(key))
}

var (
	AppName       = env("APP_NAME").String("App")
	Env           = env("ENV").String("development")
	Dev           = env("ENV").String("development") == "development"
	Prod          = env("ENV").String("production") == "production"
	CorsOrigin    = env("CORS_ORIGIN").String("(.*?)")
	Port          = env("PORT").String(":8080")
	DbUrl         = env("DB_URL").String("postgresql://postgres@localhost:5432/app?sslmode=disable")
	EmailSender   = env("EMAIL_SENDER").String("hallo@app.com")
	EmailHost     = env("EMAIL_HOST").String("smtp.app.com")
	EmailPort     = env("EMAIL_PORT").Int(587)
	EmailSecure   = env("EMAIL_SECURE").Bool(false)
	EmailUsername = env("EMAIL_USERNAME").String("hallo@app.com")
	EmailPassword = env("EMAIL_PASSWORD").String("password")
	StorageDriver = env("STORAGE_DRIVER").String("local")
)
