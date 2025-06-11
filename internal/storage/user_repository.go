package storage

import (
	"LeForum/internal/domain"
)

func SaveUserIfNotExists(email, username string) error {
	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE mail=?)", email).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		_, err := DB.Exec("INSERT INTO users (username, mail, password) VALUES (?, ?, '')", username, email)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateUser(username, email, password string) error {
	_, err := DB.Exec(
		"INSERT INTO users (username, mail, password, darkmode) VALUES (?, ?, ?, ?)",
		username, email, password, 0,
	)
	return err
}

func GetUserByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	err := DB.QueryRow(
		"SELECT id, username, mail, password FROM users WHERE mail = ?",
		email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	return user, err
}
