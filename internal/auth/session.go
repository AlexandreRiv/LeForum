package auth

import (
	"LeForum/internal/storage"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Session struct {
	ID        string
	UserEmail string
	ExpiresAt time.Time
}

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func CreateSession(w http.ResponseWriter, user LoggedUser) error {
	if user.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	sessionID := generateSessionID()
	expiresAt := time.Now().Add(5 * time.Hour)

	// Insert session into database
	result, err := storage.DB.Exec(
		"INSERT INTO sessions (id, user_email, expires_at) VALUES (?, ?, ?)",
		sessionID,
		user.Email,
		expiresAt,
	)
	if err != nil {
		log.Printf("Database insertion error: %v", err)
		return fmt.Errorf("database error: %v", err)
	}

	// Verify insertion
	rows, err := result.RowsAffected()
	if err != nil || rows != 1 {
		log.Printf("Session insertion verification failed: %v", err)
		return fmt.Errorf("session creation failed")
	}

	// Store user in manager
	manager.mu.Lock()
	manager.users[user.Email] = user
	manager.mu.Unlock()

	// Set cookie with strict settings
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Domain:   "forum.ynov.zeteox.fr",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  expiresAt,
	}
	http.SetCookie(w, cookie)

	log.Printf("Session created successfully for user: %s", user.Email)
	return nil
}

func GetSession(r *http.Request) (*Session, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, nil
		}
		log.Printf("Erreur récupération cookie: %v\n", err)
		return nil, err
	}

	log.Printf("Session ID reçu du cookie: %s", cookie.Value)

	var session Session
	var expiresAtStr string

	err = storage.DB.QueryRow(
		"SELECT id, user_email, expires_at FROM sessions WHERE id = ? AND expires_at > NOW()",
		cookie.Value,
	).Scan(&session.ID, &session.UserEmail, &expiresAtStr)

	if err == sql.ErrNoRows {
		log.Println("Aucune session trouvée dans la base pour ce cookie.")
		return nil, nil
	}
	if err != nil {
		log.Printf("Erreur DB lors de la récupération de session: %v", err)
		return nil, err
	}

	session.ExpiresAt, err = time.Parse("2006-01-02 15:04:05", expiresAtStr)
	if err != nil {
		log.Printf("Erreur parsing expires_at: %v", err)
		return nil, err
	}

	return &session, nil
}

func CleanExpiredSessions() {
	_, err := storage.DB.Exec("DELETE FROM sessions WHERE expires_at <= NOW()")
	if err != nil {
		log.Printf("Failed to clean expired sessions: %v", err)
	}
}
