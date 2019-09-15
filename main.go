package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

// OAuth2Token represents the response from Github containing the access token
type OAuth2Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func main() {
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		keys, ok := r.URL.Query()["code"]
		if !ok {
			log.Println("Authoriziation code invalid on callback")
			http.Error(w, "Authoriziation code invalid on callback", http.StatusBadRequest)
			return
		}

		ctx := context.Background()
		code := keys[0]

		keys, ok = r.URL.Query()["state"]
		if !ok {
			log.Println("State was not passed from OAuth2")
			http.Error(w, "State was not passed from OAuth2", http.StatusBadRequest)
			return
		}

		var config oauth2.Config

		state := keys[0]
		if state == "google" {
			config = GoogleConfig
		} else if state == "github" {
			config = GithubConfig
		}

		token, err := config.Exchange(ctx, code)
		if err != nil {
			log.Println(err)
			http.Error(w, "Access token missing on response", http.StatusUnprocessableEntity)
			return
		}

		oauth2Token := OAuth2Token{
			AccessToken: token.AccessToken,
			TokenType:   token.TokenType,
		}

		json, err := json.Marshal(oauth2Token)
		if err != nil {
			log.Println(err)
		}

		w.Write(json)
	})

	http.HandleFunc("/api/auth/google", func(w http.ResponseWriter, r *http.Request) {
		url := GoogleConfig.AuthCodeURL("google", oauth2.AccessTypeOnline)
		http.Redirect(w, r, url, http.StatusMovedPermanently)
	})

	http.HandleFunc("/api/auth/github", func(w http.ResponseWriter, r *http.Request) {
		url := GithubConfig.AuthCodeURL("github", oauth2.AccessTypeOnline)
		http.Redirect(w, r, url, http.StatusMovedPermanently)
	})

	http.ListenAndServe(":8080", nil)
}
