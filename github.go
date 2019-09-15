package main

import (
	"os"

	"golang.org/x/oauth2"
)

// GithubConfig is the OAuth2 config for Github
var GithubConfig = oauth2.Config{
	ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
	ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://github.com/login/oauth/authorize",
		TokenURL: "https://github.com/login/oauth/access_token",
	},
	RedirectURL: "http://localhost:8080/callback",
}
