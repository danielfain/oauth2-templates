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
	config := oauth2.Config{
		ClientID:     "c8130dc63726cbff7289",
		ClientSecret: "4a84ac7fc10c9d4f14424133e5c9531fd0d475b7",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
		RedirectURL: "http://localhost:8080/callback",
	}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		keys, ok := r.URL.Query()["code"]
		if !ok {
			log.Println("Authoriziation code invalid on callback")
			http.Error(w, "Authoriziation code invalid on callback", http.StatusBadRequest)
			return
		}

		ctx := context.Background()
		code := keys[0]

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

	http.HandleFunc("/api/auth/github", func(w http.ResponseWriter, r *http.Request) {
		url := config.AuthCodeURL("state", oauth2.AccessTypeOnline)
		http.Redirect(w, r, url, http.StatusMovedPermanently)
	})

	http.ListenAndServe(":8080", nil)

}
