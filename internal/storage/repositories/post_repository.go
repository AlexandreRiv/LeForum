package repositories

import (
	"LeForum/internal/domain"
	"database/sql"
	"time"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) CreatePost(title, content, sessionID, category string, createdAt time.Time) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Insérer le post
	postResult, err := tx.Exec(
		"INSERT INTO posts (title, content, id_user, created_at) VALUES (?, ?, (SELECT users.id FROM users INNER JOIN sessions ON users.mail = sessions.user_email WHERE sessions.id = ?), ?);",
		title,
		content,
		sessionID,
		createdAt,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Récupérer l'ID du post inséré
	postID, err := postResult.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	// Associer le post à la catégorie
	_, err = tx.Exec(
		"INSERT INTO affectation VALUES (?, (SELECT id FROM categories WHERE name = ?));",
		postID,
		category,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *PostRepository) GetPosts() ([]domain.Post, error) {
	rows, err := r.db.Query(`
		SELECT posts.id, posts.title, posts.content, users.username,
		       SUM(CASE WHEN liked_posts.liked = 1 THEN 1 ELSE 0 END) AS likes,
		       SUM(CASE WHEN liked_posts.liked = 0 THEN 1 ELSE 0 END) AS dislikes,
		       COUNT(distinct comments.id) AS comments,
		       posts.created_at
		FROM posts 
		INNER JOIN users ON posts.id_user = users.id 
		LEFT JOIN liked_posts ON liked_posts.id_post = posts.id 
		LEFT JOIN comments ON comments.id_post = posts.id 
		GROUP BY posts.id
		ORDER BY posts.created_at DESC;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []domain.Post
	for rows.Next() {
		var post domain.Post
		var createdAt string
		if err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.Username,
			&post.Likes, &post.Dislikes, &post.Comments, &createdAt); err != nil {
			return nil, err
		}
		post.CreatedAt = createdAt

		// Récupérer les catégories pour ce post
		catRows, err := r.db.Query(`
			SELECT categories.name 
			FROM categories 
			INNER JOIN affectation ON affectation.id_category = categories.id 
			WHERE affectation.id_post = ?
		`, post.Id)
		if err != nil {
			return nil, err
		}

		for catRows.Next() {
			var category string
			if err := catRows.Scan(&category); err != nil {
				catRows.Close()
				return nil, err
			}
			post.Categories = append(post.Categories, category)
		}
		catRows.Close()

		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepository) LikePost(sessionID string, postID string, likeType int) error {
	_, err := r.db.Exec(
		"INSERT INTO liked_posts VALUES ((SELECT users.id FROM users INNER JOIN sessions ON users.mail = sessions.user_email WHERE sessions.id = ?), ?, ?);",
		sessionID,
		postID,
		likeType,
	)

	if err != nil {
		_, err = r.db.Exec(
			"DELETE FROM liked_posts WHERE id_user = (SELECT users.id FROM users INNER JOIN sessions ON users.mail = sessions.user_email WHERE sessions.id = ?) AND id_post = ?;",
			sessionID,
			postID,
		)
	}

	return err
}
