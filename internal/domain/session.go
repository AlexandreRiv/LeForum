package domain

import "time"

type Session struct {
	ID        string
	UserEmail string
	ExpiresAt time.Time
}
