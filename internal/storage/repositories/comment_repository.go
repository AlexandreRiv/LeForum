package repositories

import (
	"database/sql"
	"time"
	"fmt"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) CreateComment(content, sessionID string, postID int, createdAt time.Time) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Insérer le post
	fmt.Println(content,sessionID,postID,createdAt)
	_, err = tx.Exec(
		"INSERT INTO comments (content, id_user, id_post, created_at) VALUES (?, (SELECT users.id FROM users INNER JOIN sessions ON users.mail = sessions.user_email WHERE sessions.id = ?), ?, ?);",
		content,
		sessionID,
		postID,
		createdAt,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}