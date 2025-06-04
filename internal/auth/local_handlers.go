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
	checkEmail := r.FormValue("checkEmail") == "true"

	data := map[string]interface{}{
		"RegisterEmail":    email,
		"RegisterUsername": username,
	}

	// Force le formulaire d'inscription à rester ouvert
	data["RegisterTab"] = true

	// Vérifie si c'est juste une vérification d'email
	if checkEmail {
		_, err := storage.GetUserByEmail(email)
		if err == nil {
			data["RegisterEmailError"] = "Cet email est déjà utilisé"
			h.templates.ExecuteTemplate(w, "authentification.html", data)
			return
		}
		// Email disponible, recharge le formulaire sans erreur
		h.templates.ExecuteTemplate(w, "authentification.html", data)
		return
	}

	// Processus d'inscription normal
	_, err := storage.GetUserByEmail(email)
	if err == nil {
		data["RegisterEmailError"] = "Cet email est déjà utilisé"
		h.templates.ExecuteTemplate(w, "authentification.html", data)
		return
	}

	if password != confirmPassword {
		data["RegisterPasswordError"] = "Les mots de passe ne correspondent pas"
		h.templates.ExecuteTemplate(w, "authentification.html", data)
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

func (h *Handler) handleCheckEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm-password")

	data := map[string]interface{}{
		"RegisterTab":      true,
		"RegisterEmail":    email,
		"RegisterUsername": username,
	}

	// Check if the email exists
	_, err := storage.GetUserByEmail(email)
	if err == nil {
		data["RegisterEmailError"] = "Cet email est déjà utilisé"
		h.templates.ExecuteTemplate(w, "authentification.html", data)
		return
	}

	// If we have all registration data, proceed with registration
	if username != "" && password != "" && confirmPassword != "" {
		if password != confirmPassword {
			data["RegisterPasswordError"] = "Les mots de passe ne correspondent pas"
			h.templates.ExecuteTemplate(w, "authentification.html", data)
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
		return
	}

	// Just checking email - render form with validation result
	h.templates.ExecuteTemplate(w, "authentification.html", data)
}
