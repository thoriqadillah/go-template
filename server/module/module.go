package server

import (
	"app/server"
	"app/server/module/example"
	"app/server/module/storage"
)

func init() {
	server.Register(
		example.CreateService,
		storage.CreateService,
	)
}
