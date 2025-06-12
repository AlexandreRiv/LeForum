package session

import (
	"LeForum/internal/domain"
	"log"
	"sync"
	"time"
)

type LoggedUser struct {
	Email     string
	Name      string
	LoginTime time.Time
}

type userManager struct {
	users map[string]LoggedUser
	mu    sync.RWMutex
}

var manager = &userManager{
	users: make(map[string]LoggedUser),
}

func GetUsers() []domain.LoggedUser {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	users := make([]domain.LoggedUser, 0, len(manager.users))
	for _, user := range manager.users {
		users = append(users, domain.LoggedUser{
			Email:     user.Email,
			Name:      user.Name,
			LoginTime: user.LoginTime,
		})
	}
	return users
}

func (s *Service) CleanExpiredSessions() {
	_, err := s.db.Exec("DELETE FROM sessions WHERE expires_at <= NOW()")
	if err != nil {
		log.Printf("Failed to clean expired sessions: %v", err)
	}
}
