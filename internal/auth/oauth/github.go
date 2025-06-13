package oauth

import (
	"LeForum/internal/auth/session"
	"LeForum/internal/domain"
	"LeForum/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"log"
	"net/http"
	"os"
	"time"
)

var githubOauthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
	ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	RedirectURL:  "https://forum.ynov.zeteox.fr/auth/github/callback",
	Scopes: []string{
		"user:email",
		"read:user",
	},
	Endpoint: github.Endpoint,
}

var oauthStateString = "random"

type GithubHandler struct {
	userService    *service.UserService
	sessionService *session.Service
}

func NewGithubHandler(us *service.UserService, ss *session.Service) *GithubHandler {
	return &GithubHandler{
		userService:    us,
		sessionService: ss,
	}
}

func (h *GithubHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	url := githubOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *GithubHandler) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != oauthStateString {
		http.Error(w, "État invalide", http.StatusBadRequest)
		return
	}

	token, err := githubOauthConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		log.Printf("GitHub OAuth exchange error: %v", err)
		http.Error(w, "Impossible d'échanger le code", http.StatusInternalServerError)
		return
	}

	client := githubOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		log.Printf("GitHub API error: %v", err)
		http.Error(w, "Impossible de récupérer les infos utilisateur", http.StatusInternalServerError)
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
		log.Printf("JSON decode error: %v", err)
		http.Error(w, "Erreur de décodage", http.StatusInternalServerError)
		return
	}

	// Get email if not provided in main profile
	if githubUser.Email == "" {
		emails, err := getGithubEmails(client)
		if err != nil {
			log.Printf("GitHub emails error: %v", err)
			http.Error(w, "Impossible de récupérer l'email", http.StatusInternalServerError)
			return
		}
		if len(emails) > 0 {
			githubUser.Email = emails[0]
		} else {
			log.Printf("No email found for GitHub user")
			http.Error(w, "Aucun email trouvé", http.StatusInternalServerError)
			return
		}
	}

	// Use login name if no name provided
	if githubUser.Name == "" {
		githubUser.Name = githubUser.Login
	}

	// Verify we have required data
	if githubUser.Email == "" {
		log.Printf("No email available for GitHub user")
		http.Error(w, "Email requis", http.StatusInternalServerError)
		return
	}

	// Save user first
	if err := h.userService.SaveUserIfNotExists(githubUser.Email, githubUser.Name); err != nil {
		log.Printf("Error saving user: %v", err)
		http.Error(w, "Erreur sauvegarde utilisateur", http.StatusInternalServerError)
		return
	}

	user := domain.LoggedUser{
		Email:     githubUser.Email,
		Name:      githubUser.Name,
		LoginTime: time.Now(),
	}

	if err := h.sessionService.CreateSession(w, user); err != nil {
		log.Printf("Session creation error: %v", err)
		http.Error(w, "Erreur création session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
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
