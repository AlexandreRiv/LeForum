package auth

import (
	"LeForum/internal/api"
	"html/template"
	"log"
	"net/http"
)

type Handler struct {
	templates *template.Template
}

func NewHandler() *Handler {
	// Chargement explicite de tous les templates nécessaires
	tmpl := template.New("")

	// Charger d'abord les templates de base
	baseTemplates, err := tmpl.ParseGlob("web/templates/*.html")
	if err != nil {
		log.Fatalf("Erreur de chargement des templates de base: %v", err)
	}

	// Charger les templates de composants
	_, err = baseTemplates.ParseGlob("web/templates/components/*.html")
	if err != nil {
		log.Fatalf("Erreur de chargement des templates de composants: %v", err)
	}

	return &Handler{templates: baseTemplates}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := GetSession(r)
	if err != nil {
		log.Println("Erreur session:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if session != nil {
		log.Println("Session déjà active, redirection")
		http.Redirect(w, r, "/users", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		switch r.FormValue("action") {
		case "register":
			h.handleRegister(w, r)
		case "login":
			h.handleLogin(w, r)
		default:
			http.Error(w, "Invalid action", http.StatusBadRequest)
		}
		return
	}

	// Get tab from URL parameter, default to login
	tab := r.URL.Query().Get("tab")
	data := map[string]interface{}{
		"RegisterTab": tab == "register",
		"DarkMode":    getDarkModeFromCookie(r),
		"CurrentPage": "auth",
	}

	err = h.templates.ExecuteTemplate(w, "authentification", data)
	if err != nil {
		log.Printf("Erreur rendering template: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("/auth", h)
	mux.HandleFunc("/auth/check-email", h.handleCheckEmail)
	mux.HandleFunc("/users", h.UserPageHandler)
	mux.HandleFunc("/auth/google", GoogleLoginHandler)
	mux.HandleFunc("/auth/google/callback", GoogleCallbackHandler)
	mux.HandleFunc("/auth/github", GithubLoginHandler)
	mux.HandleFunc("/auth/github/callback", GithubCallbackHandler)
	mux.HandleFunc("/admin", AdminHandler)
	mux.HandleFunc("/logout", h.LogoutHandler)
}

func (h *Handler) HandleAuth(w http.ResponseWriter, r *http.Request) {
	data := api.PageData{
		DarkMode:    getDarkModeFromCookie(r),
		CurrentPage: "auth",
	}

	if err := h.templates.ExecuteTemplate(w, "authentification", data); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func getDarkModeFromCookie(r *http.Request) bool {
	cookie, err := r.Cookie("darkMode")
	if err == nil && cookie.Value == "true" {
		return true
	}
	return false
}
