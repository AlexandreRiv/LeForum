package domain

import "time"

type Session struct {
	ID        string    `json:"id"`
	UserEmail string    `json:"user_email"`
	ExpiresAt time.Time `json:"expires_at"`
}
