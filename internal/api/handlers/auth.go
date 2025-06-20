package handlers

import (
	"LeForum/internal/auth/session"
	"LeForum/internal/domain"
	"LeForum/internal/service"
	"log"
	"net/http"
	"regexp"
	"time"
)

type AuthHandler struct {
	userService     *service.UserService
	notificationService *service.NotificationService
	sessionService  *session.Service
	templateService *TemplateService
}

func NewAuthHandler(us *service.UserService, ns *service.NotificationService, ss *session.Service, ts *TemplateService) *AuthHandler {
	return &AuthHandler{
		userService:     us,
		notificationService: ns,
		sessionService:  ss,
		templateService: ts,
	}
}

func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/auth", h.HandleAuth)
	mux.HandleFunc("/auth/check-email", h.handleCheckEmail)
	mux.HandleFunc("/users", h.UserPageHandler)
	mux.HandleFunc("/logout", h.LogoutHandler)
}

func (h *AuthHandler) HandleAuth(w http.ResponseWriter, r *http.Request) {
	session, err := h.sessionService.GetSession(r)
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

	err = h.templateService.RenderTemplate(w, "authentification.html", data)
	if err != nil {
		log.Printf("Erreur rendering template: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *AuthHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
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
		_, err := h.userService.GetUserByEmail(email)
		if err == nil {
			data["RegisterEmailError"] = "Cet email est déjà utilisé"
			h.templateService.RenderTemplate(w, "authentification.html", data)
			return
		}
		// Email disponible, recharge le formulaire sans erreur
		h.templateService.RenderTemplate(w, "authentification.html", data)
		return
	}

	// Processus d'inscription normal
	_, err := h.userService.GetUserByEmail(email)
	if err == nil {
		data["RegisterEmailError"] = "Cet email est déjà utilisé"
		h.templateService.RenderTemplate(w, "authentification.html", data)
		return
	}

	if password != confirmPassword {
		data["RegisterPasswordError"] = "Les mots de passe ne correspondent pas"
		h.templateService.RenderTemplate(w, "authentification.html", data)
		return
	}

	err = h.userService.CreateUser(username, email, password)
	if err != nil {
		log.Printf("Erreur lors de la création de l'utilisateur: %v", err)
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	user := domain.LoggedUser{
		Email:     email,
		Name:      username,
		LoginTime: time.Now(),
	}
	h.sessionService.CreateSession(w, user)
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func (h *AuthHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	data := make(map[string]interface{})

	user, err := h.userService.GetUserByEmail(email)
	if err != nil {
		log.Printf("Login failed for email %s: %v", email, err)
		data["EmailError"] = "Email invalide"
		h.templateService.RenderTemplate(w, "authentification.html", data)
		return
	}

	if !h.userService.VerifyPassword(password, user.Password) {
		log.Printf("Invalid password for email %s", email)
		data["PasswordError"] = "Mot de passe incorrect"
		h.templateService.RenderTemplate(w, "authentification.html", data)
		return
	}

	loggedUser := domain.LoggedUser{
		Email:     user.Email,
		Name:      user.Username,
		LoginTime: time.Now(),
	}

	if err := h.sessionService.CreateSession(w, loggedUser); err != nil {
		log.Printf("Failed to create session for %s: %v", email, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) handleCheckEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	username := r.FormValue("username")
	data := map[string]interface{}{
		"RegisterTab":      true,
		"RegisterEmail":    email,
		"RegisterUsername": username,
	}

	// Regex to validate email format
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	isValidEmail := regexp.MustCompile(emailRegex).MatchString(email)

	if !isValidEmail {
		data["RegisterEmailError"] = "L'email n'est pas valide"
		h.templateService.RenderTemplate(w, "authentification.html", data)
		return
	}

	// Check if the email exists
	_, err := h.userService.GetUserByEmail(email)
	if err == nil {
		data["RegisterEmailError"] = "Cet email est déjà utilisé"
	}

	// Render the template with the register tab active
	h.templateService.RenderTemplate(w, "authentification.html", data)
}

func (h *AuthHandler) UserPageHandler(w http.ResponseWriter, r *http.Request) {
	session, err := h.sessionService.GetSession(r)
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

	user, err := h.userService.GetUserByEmail(session.UserEmail)
	if err != nil {
		http.Error(w, "Failed to fetch user data", http.StatusInternalServerError)
		return
	}

	loggedUser := &domain.LoggedUser{
		Email: user.Email,
		Name:  user.Username,
	}

	notifs, err := h.notificationService.GetNotifications(session.ID)

	data := struct {
		Name        string
		Email       string
		Expiry      time.Time
		DarkMode    bool
		CurrentPage string
		User        *domain.LoggedUser
		Likes       int
		PostNumber  int
		ResponseNb  int
		Notifications []domain.Notification
		NotificationNb int
	}{
		Name:        user.Username,
		Email:       user.Email,
		Expiry:      session.ExpiresAt,
		DarkMode:    getDarkModeFromCookie(r),
		CurrentPage: "profile",
		User:        loggedUser,
		Notifications: notifs,
		NotificationNb: len(notifs),
	}

	// Logique pour récupérer les statistiques utilisateur
	// À implémenter dans un service approprié

	PostNb, RespNb, Likes, _ := h.userService.GetUserStats(loggedUser.Email)
	data.Likes = Likes
	data.PostNumber = PostNb
	data.ResponseNb = RespNb

	err = h.templateService.RenderTemplate(w, "user.html", data)
	if err != nil {
		log.Printf("Erreur rendering user template: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer le cookie de session
	cookie, err := r.Cookie("session_id")
	if err == nil {
		// Utiliser une requête SQL directement car le service session n'expose pas de méthode DeleteSession
		_, err = h.sessionService.DeleteSession(cookie.Value)
		if err != nil {
			log.Printf("Erreur lors de la suppression de la session: %v", err)
			// Continuer malgré l'erreur pour expirer le cookie côté client
		}
	}

	// Créer un cookie expiré avec les mêmes propriétés que celui utilisé à la création
	expiredCookie := &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		MaxAge:   -1, // Expire immédiatement
	}

	http.SetCookie(w, expiredCookie)
	http.Redirect(w, r, "/auth", http.StatusSeeOther)
}
