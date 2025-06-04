package auth

import (
	"html/template"
	"log"
	"net/http"
	"time"
	"LeForum/internal/storage"
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

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
    email := r.FormValue("email")
    username := r.FormValue("username")
    password := r.FormValue("password")
    confirmPassword := r.FormValue("confirm-password")

    if password != confirmPassword {
        http.Error(w, "Passwords don't match", http.StatusBadRequest)
        return
    }

    hashedPassword, err := HashPassword(password)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    err = storage.CreateUser(username, email, hashedPassword)
    if err != nil {
	log.Printf("Erreur lors de la création de l'utilisateur: %v", err)
        http.Error(w, "Could not create user", http.StatusInternalServerError)
        return
    }

    // Create session and redirect
    user := LoggedUser{
        Email: email,
        Name: username,
        LoginTime: time.Now(),
    }
    CreateSession(w, user)
    http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
    email := r.FormValue("email")
    password := r.FormValue("password")

    user, err := storage.GetUserByEmail(email)
    if err != nil {
        log.Printf("Login failed for email %s: %v", email, err)
        w.WriteHeader(http.StatusUnauthorized)
        w.Write([]byte("Invalid credentials"))
        return
    }

    if !CheckPassword(password, user.Password) {
        log.Printf("Invalid password for email %s", email)
        w.WriteHeader(http.StatusUnauthorized)
        w.Write([]byte("Invalid credentials"))
        return
    }

    loggedUser := LoggedUser{
        Email: user.Email,
        Name:  user.Username,
        LoginTime: time.Now(),
    }

    if err := CreateSession(w, loggedUser); err != nil {
        log.Printf("Failed to create session for %s: %v", email, err)
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Internal server error"))
        return
    }

    http.Redirect(w, r, "/users", http.StatusSeeOther)
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
