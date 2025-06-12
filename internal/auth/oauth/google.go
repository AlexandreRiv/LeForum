package oauth

import (
	"LeForum/internal/auth/session"
	"LeForum/internal/domain"
	"LeForum/internal/service"
	"context"
	"encoding/json"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"net/http"
	"os"
	"time"
)

// Configuration OAuth pour Google
var googleOauthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL:  "https://forum.ynov.zeteox.fr/auth/google/callback",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}

type GoogleHandler struct {
	userService    *service.UserService
	sessionService *session.Service
}

func NewGoogleHandler(us *service.UserService, ss *session.Service) *GoogleHandler {
	return &GoogleHandler{
		userService:    us,
		sessionService: ss,
	}
}

func (h *GoogleHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *GoogleHandler) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != oauthStateString {
		http.Error(w, "État invalide", http.StatusBadRequest)
		return
	}

	token, err := googleOauthConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		log.Printf("Erreur échange code: %v", err)
		http.Error(w, "Impossible d'échanger le code", http.StatusInternalServerError)
		return
	}

	client := googleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Printf("Erreur récupération infos: %v", err)
		http.Error(w, "Impossible de récupérer les infos utilisateur", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		log.Printf("Erreur décodage JSON: %v", err)
		http.Error(w, "Erreur de décodage", http.StatusInternalServerError)
		return
	}

	email := userInfo["email"].(string)
	name := userInfo["name"].(string)

	// Save user if not exists
	err = h.userService.SaveUserIfNotExists(email, name)
	if err != nil {
		log.Printf("Erreur sauvegarde utilisateur: %v", err)
		http.Error(w, "Erreur sauvegarde utilisateur", http.StatusInternalServerError)
		return
	}

	user := domain.LoggedUser{
		Email:     email,
		Name:      name,
		LoginTime: time.Now(),
	}

	// Create session with error handling
	err = h.sessionService.CreateSession(w, user)
	if err != nil {
		log.Printf("Erreur création session: %v", err)
		http.Error(w, "Erreur création session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
