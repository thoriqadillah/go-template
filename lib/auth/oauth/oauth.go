package oauth

import (
	"context"
	"log"
)

type User struct {
	Email string `json:"email"`
	Image string `json:"picture"`
	Name  string `json:"name"`
}

type OAuth interface {
	Validate(ctx context.Context, token string) (*User, error)
}

var providers = make(map[string]OAuth)

func register(name string, impl OAuth) {
	providers[name] = impl
}

func Create(provider string) OAuth {
	oauth, ok := providers[provider]
	if !ok {
		log.Fatalf("Oauth with %s provider is not implemented", provider)
		return nil
	}

	return oauth
}
