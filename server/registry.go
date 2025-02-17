package server

import (
	"app/server/example"
	"app/server/storage"
)

func init() {
	Register(
		example.CreateService(),
		storage.CreateService(),
	)
}
