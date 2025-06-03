package auth

import (
	"html/template"
	"log"
	"net/http"
)

type Handler struct {
	templates *template.Template
}

func NewHandler() *Handler {
	return &Handler{
		templates: template.Must(template.ParseFiles("web/templates/authentification.html")),
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.templates.ExecuteTemplate(w, "authentification.html", nil)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("/auth", h)
	mux.HandleFunc("/auth/google", googleLoginHandler)
	mux.HandleFunc("/auth/google/callback", googleCallbackHandler)
	mux.HandleFunc("/auth/github", githubLoginHandler)
	mux.HandleFunc("/auth/github/callback", githubCallbackHandler)
	mux.HandleFunc("/admin", adminHandler)
}
