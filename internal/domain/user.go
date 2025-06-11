package domain

import "time"

type User struct {
	ID       int
	Username string
	Email    string
	Password string
	DarkMode bool
}

type LoggedUser struct {
	Email     string
	Name      string
	LoginTime time.Time
}
