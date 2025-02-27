package oauth

import "context"

type googleOAuth struct{}

func createGoogleOAuth() OAuth {
	return &googleOAuth{}
}

func (g *googleOAuth) Validate(ctx context.Context, token string) (*User, error) {
	panic("TODO: implement google oauth")
}

func init() {
	register("google", createGoogleOAuth())
}
