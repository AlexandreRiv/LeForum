package domain

import (
	"LeForum/internal/api/middleware"
	"time"
)

type LoggedUser struct {
	Email     string              `json:"email"`
	Name      string              `json:"name"`
	LoginTime time.Time           `json:"login_time"`
	Role      middleware.RoleType `json:"role"`
}
