package domain

import "LeForum/internal/api/middleware"

type User struct {
	ID       int
	Username string
	Email    string
	Password string
	DarkMode bool
	Role     middleware.RoleType
}
