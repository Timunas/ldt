package google

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const GoogleUserInfoEndpoint = "https://www.googleapis.com/oauth2/v3/userinfo"

func NewConfig(clientID string, clientSecret string, callbackURL string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  callbackURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
