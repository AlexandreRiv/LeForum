package repositories

import (
	"LeForum/internal/domain"
	"database/sql"
	"log"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) SaveUserIfNotExists(email, username string) error {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE mail = ?", email).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		_, err = r.db.Exec("INSERT INTO users (username, mail) VALUES (?, ?)", username, email)
	}
	return err
}

func (r *UserRepository) CreateUser(username, email, hashedPassword string) error {
	_, err := r.db.Exec("INSERT INTO users (username, mail, password) VALUES (?, ?, ?)", username, email, hashedPassword)
	return err
}

func (r *UserRepository) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow("SELECT id, username, mail, password FROM users WHERE mail = ?", email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		log.Printf("Erreur lors de la récupération de l'utilisateur: %v", err)
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserStats(email string) (postCount, responseCount, likeCount int, err error) {
	err = r.db.QueryRow(
		"SELECT COUNT(*) FROM posts WHERE posts.id_user = (SELECT id FROM users WHERE mail = ?)",
		email,
	).Scan(&postCount)
	if err != nil {
		return 0, 0, 0, err
	}

	err = r.db.QueryRow(
		"SELECT COUNT(*) FROM comments WHERE comments.id_user = (SELECT id FROM users WHERE mail = ?)",
		email,
	).Scan(&responseCount)
	if err != nil {
		return postCount, 0, 0, err
	}

	err = r.db.QueryRow(
		"SELECT COUNT(*) FROM liked_posts INNER JOIN posts ON liked_posts.id_post = posts.id "+
			"WHERE liked_posts.liked = 1 AND posts.id_user = (SELECT id FROM users WHERE mail = ?)",
		email,
	).Scan(&likeCount)

	return postCount, responseCount, likeCount, err
}
