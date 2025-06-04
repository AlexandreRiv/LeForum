package auth

import (
	"LeForum/internal/storage"
	"log"
	"net/http"
	"time"
)

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

	manager.mu.RLock()
	user, exists := manager.users[session.UserEmail]
	manager.mu.RUnlock()

	if !exists {
		user = LoggedUser{
			Email: session.UserEmail,
			Name:  session.UserEmail,
		}
	}

	data := struct {
		Name        string
		Email       string
		Expiry      time.Time
		DarkMode    bool
		CurrentPage string
	}{
		Name:        user.Name,
		Email:       user.Email,
		Expiry:      session.ExpiresAt,
		DarkMode:    getDarkModeFromCookie(r),
		CurrentPage: "profile",
	}

	err = h.templates.ExecuteTemplate(w, "user.html", data)
	if err != nil {
		log.Printf("Erreur rendering user template: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		_, err = storage.DB.Exec("DELETE FROM sessions WHERE id = ?", cookie.Value)
		if err != nil {
			log.Printf("Error deleting session: %v", err)
		}
	}

	expiredCookie := &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Domain:   "forum.ynov.zeteox.fr",
		Expires:  time.Now().Add(-24 * time.Hour),
	}
	http.SetCookie(w, expiredCookie)

	http.Redirect(w, r, "/auth", http.StatusSeeOther)
}
