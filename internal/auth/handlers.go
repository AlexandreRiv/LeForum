package auth

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	templates *template.Template
}

func NewHandler() *Handler {
    tmpl, err := template.ParseGlob("web/templates/*.html")
    if err != nil {
        log.Fatalf("Erreur de chargement des templates: %v", err)
    }

    for _, t := range tmpl.Templates() {
        log.Printf("Template loaded: %s", t.Name())
    }

    return &Handler{templates: tmpl}
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

    log.Println("Rendering authentification.html template")
    err = h.templates.ExecuteTemplate(w, "authentification.html", nil)
    if err != nil {
        log.Printf("Erreur rendering template: %v\n", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}

func (h *Handler) UserPageHandler(w http.ResponseWriter, r *http.Request) {
    session, err := GetSession(r)
    if err != nil {
        log.Printf("Erreur récupération session: %v\n", err)
        http.Redirect(w, r, "/auth", http.StatusSeeOther)
        return
    }
    if session == nil {
        log.Println("Session vide, redirection vers /auth")
        http.Redirect(w, r, "/auth", http.StatusSeeOther)
        return
    }

    log.Printf("Session OK: email = %s, expiry = %v\n", session.UserEmail, session.ExpiresAt)

    data := struct {
        Email  string
        Expiry time.Time
    }{
        Email:  session.UserEmail,
        Expiry: session.ExpiresAt,
    }

    err = h.templates.ExecuteTemplate(w, "user.html", data)
    if err != nil {
        log.Printf("Erreur rendering user template: %v\n", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("/auth", h)
	mux.HandleFunc("/users", h.UserPageHandler)
	mux.HandleFunc("/auth/google", GoogleLoginHandler)
	mux.HandleFunc("/auth/google/callback", GoogleCallbackHandler)
	mux.HandleFunc("/auth/github", GithubLoginHandler)
	mux.HandleFunc("/auth/github/callback", GithubCallbackHandler)
	mux.HandleFunc("/admin", AdminHandler)
}
