package domain

import "time"

type LoggedUser struct {
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	LoginTime time.Time `json:"login_time"`
}
