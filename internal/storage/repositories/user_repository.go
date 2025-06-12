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

func (r *UserRepository) UpdateUserRole(userID int, role string) error {
	_, err := r.db.Exec("UPDATE users SET user_role = ? WHERE id = ?", role, userID)
	return err
}

func (r *UserRepository) GetAllUsers() ([]*domain.User, error) {
	rows, err := r.db.Query("SELECT id, username, mail, password, user_role FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*domain.User{}
	for rows.Next() {
		user := &domain.User{}
		var roleStr string
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &roleStr)
		if err != nil {
			return nil, err
		}
		user.Role = domain.RoleType(roleStr)
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	var roleStr string
	err := r.db.QueryRow("SELECT id, username, mail, password, user_role FROM users WHERE mail = ?", email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &roleStr)
	if err != nil {
		log.Printf("Erreur lors de la récupération de l'utilisateur: %v", err)
		return nil, err
	}
	user.Role = domain.RoleType(roleStr)
	return &user, nil
}

func (r *UserRepository) GetUserByID(id int) (*domain.User, error) {
	var user domain.User
	var roleStr string
	err := r.db.QueryRow("SELECT id, username, mail, password, user_role FROM users WHERE id = ?", id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &roleStr)
	if err != nil {
		return nil, err
	}
	user.Role = domain.RoleType(roleStr)
	return &user, nil
}
