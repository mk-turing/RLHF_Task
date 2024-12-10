package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	redirectURL = "http://localhost:8080/callback"
	scope       = "https://www.googleapis.com/auth/userinfo.profile https://www.googleapis.com/auth/userinfo.email"
)

var config = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL:  redirectURL,
	Scopes:       []string{scope},
	Endpoint:     google.Endpoint,
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	authURL := config.AuthCodeURL(r.URL.Query().Get("state"))
	http.Redirect(w, r, authURL, http.StatusFound)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	token, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
		return
	}

	// Retrieve user information
	client := config.Client(oauth2.NoContext, token)
	userInfoURL := "https://www.googleapis.com/oauth2/v1/userinfo"
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var user struct {
		ID             string `json:"id"`
		Email          string `json:"email"`
		Name           string `json:"name"`
		Verified_email bool   `json:"verified_email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Logged in as %s (%s)\n", user.Name, user.Email)
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/callback", callbackHandler)

	fmt.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
