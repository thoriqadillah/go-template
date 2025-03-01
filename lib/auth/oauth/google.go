package oauth

import (
	"app/env"
	"context"
	"encoding/json"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type userInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Locale        string `json:"locale"`
}

var googleConfig = &oauth2.Config{
	ClientID:     env.GOOGLE_OAUTH_CLIENT_ID,
	ClientSecret: env.GOOGLE_OAUTH_SECRET,
	RedirectURL:  "",         // TODO
	Scopes:       []string{}, // TODO
	Endpoint:     google.Endpoint,
}

type googleOAuth struct{}

func createGoogleOAuth() OAuth {
	return &googleOAuth{}
}

func (g *googleOAuth) Validate(ctx context.Context, token string) (*User, error) {
	tok := oauth2.Token{AccessToken: token}
	tokenSource := googleConfig.TokenSource(ctx, &tok)
	client := oauth2.NewClient(ctx, tokenSource)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo userInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	user := &User{
		Email: userInfo.Email,
		Image: userInfo.Picture,
		Name:  userInfo.Name,
	}

	return user, nil
}

func init() {
	register("google", createGoogleOAuth())
}
