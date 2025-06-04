package auth

import (
	"LeForum/internal/storage"
	"log"
	"net/http"
	"time"
)

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm-password")

	data := make(map[string]interface{})

	// Vérifier si l'email existe déjà
	_, err := storage.GetUserByEmail(email)
	if err == nil {
		// L'email existe déjà
		data["RegisterEmailError"] = "Cet email est déjà utilisé"
		h.templates.ExecuteTemplate(w, "authentification.html", data)
		return
	}

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

	user := LoggedUser{
		Email:     email,
		Name:      username,
		LoginTime: time.Now(),
	}
	CreateSession(w, user)
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	data := make(map[string]interface{})

	user, err := storage.GetUserByEmail(email)
	if err != nil {
		log.Printf("Login failed for email %s: %v", email, err)
		data["EmailError"] = "Email invalide"
		h.templates.ExecuteTemplate(w, "authentification.html", data)
		return
	}

	if !CheckPassword(password, user.Password) {
		log.Printf("Invalid password for email %s", email)
		data["PasswordError"] = "Mot de passe incorrect"
		h.templates.ExecuteTemplate(w, "authentification.html", data)
		return
	}

	loggedUser := LoggedUser{
		Email:     user.Email,
		Name:      user.Username,
		LoginTime: time.Now(),
	}

	if err := CreateSession(w, loggedUser); err != nil {
		log.Printf("Failed to create session for %s: %v", email, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
