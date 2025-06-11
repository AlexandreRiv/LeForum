package session

import (
	"LeForum/internal/domain"
	"crypto/rand"
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateSession(w http.ResponseWriter, user domain.LoggedUser) error {
	sessionID := s.GenerateSessionID()
	expiryTime := time.Now().Add(24 * time.Hour)

	_, err := s.db.Exec(
		"INSERT INTO sessions (id, user_email, expires_at) VALUES (?, ?, ?)",
		sessionID, user.Email, expiryTime,
	)
	if err != nil {
		return err
	}

	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  expiryTime,
	}
	http.SetCookie(w, &cookie)

	// Utiliser la structure LoggedUser locale du package session
	manager.mu.Lock()
	manager.users[user.Email] = LoggedUser{
		Email:     user.Email,
		Name:      user.Name,
		LoginTime: user.LoginTime,
	}
	manager.mu.Unlock()

	return nil
}

func (s *Service) GetSession(r *http.Request) (*domain.Session, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, nil // Pas d'erreur, juste pas de session
	}

	var session domain.Session
	var expiresAtStr string

	err = s.db.QueryRow(
		"SELECT id, user_email, expires_at FROM sessions WHERE id = ?",
		cookie.Value,
	).Scan(&session.ID, &session.UserEmail, &expiresAtStr)
	if err != nil {
		return nil, err
	}

	// Conversion de la chaîne en time.Time
	expiresAt, err := time.Parse("2006-01-02 15:04:05", expiresAtStr)
	if err != nil {
		return nil, fmt.Errorf("impossible de parser la date d'expiration: %w", err)
	}
	session.ExpiresAt = expiresAt

	return &session, nil
}

func (s *Service) GetCurrentUser(r *http.Request) (*domain.LoggedUser, error) {
	session, err := s.GetSession(r)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, nil
	}

	manager.mu.RLock()
	defer manager.mu.RUnlock()

	// Vérifier si l'utilisateur existe dans le manager
	if user, exists := manager.users[session.UserEmail]; exists {
		return &domain.LoggedUser{
			Email:     user.Email,
			Name:      user.Name,
			LoginTime: user.LoginTime,
		}, nil
	}

	// Si l'utilisateur n'existe pas dans le manager, essayer de le récupérer depuis la base de données
	var username string
	err = s.db.QueryRow("SELECT username FROM users WHERE mail = ?", session.UserEmail).Scan(&username)
	if err != nil {
		return nil, err
	}

	// Créer un nouvel utilisateur avec les informations disponibles
	user := &domain.LoggedUser{
		Email:     session.UserEmail,
		Name:      username,
		LoginTime: time.Now(), // Pas la vraie date de login, mais une approximation
	}

	return user, nil
}

func (s *Service) GenerateSessionID() string {
	// Générer un ID de session aléatoire
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func (s *Service) DeleteSession(sessionID string) (sql.Result, error) {
	return s.db.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
}
