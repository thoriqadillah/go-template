package env

import (
	"app/common"
	"os"
)

func env(key string) common.Parser {
	return common.Parse(os.Getenv(key))
}

var (
	CorsOrigin = env("CORS_ORIGIN").String("(.*?)")
	Port       = env("PORT").String(":8080")
)
