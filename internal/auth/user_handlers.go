package auth

import (
	"LeForum/internal/domain"
	"LeForum/internal/storage"
	"log"
	"net/http"
	"os"
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
		User        *LoggedUser
		Likes       int
		PostNumber  int
		ResponseNb  int
		Posts       []domain.Post
	}{
		Name:        user.Name,
		Email:       user.Email,
		Expiry:      session.ExpiresAt,
		DarkMode:    getDarkModeFromCookie(r),
		CurrentPage: "profile",
		User:        &user,
	}

	SQLPostNbReq := "SELECT COUNT(*) FROM posts WHERE posts.id_user = (SELECT id FROM users WHERE mail = ?);"
	SQLRespNbReq := "SELECT COUNT(*) FROM comments WHERE comments.id_user = (SELECT id FROM users WHERE mail = ?);"
	SQLLikeNbReq := "SELECT COUNT(*) FROM liked_posts INNER JOIN posts ON liked_posts.id_post = posts.id WHERE liked_posts.liked = 1 AND posts.id_user = (SELECT id FROM users WHERE mail = ?);"

	err = storage.DB.QueryRow(SQLPostNbReq, data.Email).Scan(&data.PostNumber)
	if err != nil {
		http.Error(w, "Failed to fetch post number", http.StatusInternalServerError)
		return
	}
	err = storage.DB.QueryRow(SQLRespNbReq, data.Email).Scan(&data.ResponseNb)
	if err != nil {
		http.Error(w, "Failed to fetch responses number", http.StatusInternalServerError)
		return
	}
	err = storage.DB.QueryRow(SQLLikeNbReq, data.Email).Scan(&data.Likes)
	if err != nil {
		http.Error(w, "Failed to fetch responses number", http.StatusInternalServerError)
		return
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

	// Set cookie parameters based on environment
	secure := true
	domain := "forum.ynov.zeteox.fr"

	// En développement, les cookies ne nécessitent pas ces restrictions
	if os.Getenv("ENVIRONMENT") != "production" {
		secure = false
		domain = "" // Ne pas définir de domaine en développement
	}

	expiredCookie := &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		Domain:   domain,
		Expires:  time.Now().Add(-24 * time.Hour),
	}
	http.SetCookie(w, expiredCookie)

	http.Redirect(w, r, "/auth", http.StatusSeeOther)
}
