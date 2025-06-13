package repositories

import (
	"LeForum/internal/domain"
	"database/sql"
	"time"
	"fmt"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) CreateNotification(postID int, content string, createdAt time.Time) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Insérer le post
	_, err = tx.Exec(
		"INSERT INTO notifications (id_creator, id_post, content, created_at) VALUES ((SELECT users.id FROM users INNER JOIN posts ON users.id = posts.id_user WHERE posts.id = ?), ?, ?, ?);",
		postID,
		postID,
		content,
		createdAt,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *NotificationRepository) GetNotifications(sessionID string) ([]domain.Notification, error) {
	Request := "SELECT id,content,created_at FROM notifications WHERE id_creator = (SELECT users.id FROM users INNER JOIN sessions ON users.mail = sessions.user_email WHERE sessions.id = ?);"

	rows, err := r.db.Query(Request, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifs []domain.Notification
	for rows.Next() {
		var notif domain.Notification
		var createdAt string
		if err := rows.Scan(&notif.Id, &notif.Content, &createdAt); err != nil {
			fmt.Println(err)
			return nil, err
		}
		notif.CreatedAt = createdAt

		notifs = append(notifs, notif)
	}

	return notifs, nil
}