package server

import "app/server/example"

func init() {
	Register(example.CreateService())
}
