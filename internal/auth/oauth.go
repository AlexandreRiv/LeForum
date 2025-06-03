package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"html/template"
	"net/http"
	"sync"
	"time"
)

// Configuration OAuth pour Google
var googleOauthConfig = &oauth2.Config{
	ClientID:     "47465948270-8u1i4r4jlfhsnjcqusjfkk2cjsu57nef.apps.googleusercontent.com",
	ClientSecret: "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",
	RedirectURL:  "http://localhost:3002/auth/google/callback",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}

// Configuration OAuth pour GitHub
var githubOauthConfig = &oauth2.Config{
	ClientID:     "Ov23liSzicwlo2PrVNtL",
	ClientSecret: "82c9b976c0f73ea80b31e4a44f7265bed983a75d",
	RedirectURL:  "http://localhost:3002/auth/github/callback",
	Scopes: []string{
		"user:email",
		"read:user",
	},
	Endpoint: github.Endpoint,
}

var oauthStateString = "random"

// LoggedUser représente un utilisateur connecté
type LoggedUser struct {
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Picture   string    `json:"picture"`
	LoginTime time.Time `json:"login_time"`
}

type userManager struct {
	users map[string]LoggedUser
	mu    sync.RWMutex
}

var manager = &userManager{
	users: make(map[string]LoggedUser),
}

func GetUsers() []LoggedUser {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	users := make([]LoggedUser, 0, len(manager.users))
	for _, user := range manager.users {
		users = append(users, user)
	}
	return users
}

func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func googleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != oauthStateString {
		http.Error(w, "État invalide", http.StatusBadRequest)
		return
	}

	token, err := googleOauthConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		http.Error(w, "Impossible d'échanger le code : "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := googleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Impossible de récupérer les infos utilisateur : "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Erreur de décodage: "+err.Error(), http.StatusInternalServerError)
		return
	}

	manager.mu.Lock()
	manager.users[userInfo["email"].(string)] = LoggedUser{
		Email:     userInfo["email"].(string),
		Name:      userInfo["name"].(string),
		Picture:   userInfo["picture"].(string),
		LoginTime: time.Now(),
	}
	manager.mu.Unlock()

	fmt.Printf("Nouvel utilisateur Google connecté : %s\n", userInfo["email"])
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func githubLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := githubOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func githubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != oauthStateString {
		http.Error(w, "État invalide", http.StatusBadRequest)
		return
	}

	token, err := githubOauthConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		http.Error(w, "Impossible d'échanger le code : "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := githubOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		http.Error(w, "Impossible de récupérer les infos utilisateur : "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var githubUser struct {
		Email     string `json:"email"`
		Name      string `json:"name"`
		Login     string `json:"login"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		http.Error(w, "Erreur de décodage: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if githubUser.Email == "" {
		emails, err := getGithubEmails(client)
		if err == nil && len(emails) > 0 {
			githubUser.Email = emails[0]
		}
	}

	if githubUser.Name == "" {
		githubUser.Name = githubUser.Login
	}

	manager.mu.Lock()
	manager.users[githubUser.Email] = LoggedUser{
		Email:     githubUser.Email,
		Name:      githubUser.Name,
		Picture:   githubUser.AvatarURL,
		LoginTime: time.Now(),
	}
	manager.mu.Unlock()

	fmt.Printf("Nouvel utilisateur GitHub connecté : %s\n", githubUser.Email)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func getGithubEmails(client *http.Client) ([]string, error) {
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return nil, err
	}

	// Retourner le premier email vérifié et primaire
	for _, email := range emails {
		if email.Primary && email.Verified {
			return []string{email.Email}, nil
		}
	}

	return nil, fmt.Errorf("no primary verified email found")
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	users := GetUsers()
	tmpl := template.Must(template.ParseFiles("web/templates/admin.html"))
	err := tmpl.Execute(w, map[string]interface{}{
		"Users": users,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
