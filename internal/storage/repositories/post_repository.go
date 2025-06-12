package repositories

import (
	"LeForum/internal/domain"
	"database/sql"
	"time"
	"encoding/base64"
	"html/template"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) CreatePost(title, content, sessionID, category string, image []byte, createdAt time.Time) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Insérer le post
	postResult, err := tx.Exec(
		"INSERT INTO posts (title, content, id_user, image, created_at) VALUES (?, ?, (SELECT users.id FROM users INNER JOIN sessions ON users.mail = sessions.user_email WHERE sessions.id = ?), ?, ?);",
		title,
		content,
		sessionID,
		image,
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

func (r *PostRepository) GetPosts(order, search string) ([]domain.Post, error) {
	var Request string 
	BaseRequest := `
		SELECT 
			posts.id, 
			posts.title, 
			posts.content, 
			users.username,
			COALESCE(like_stats.likes, 0) AS likes,
			COALESCE(like_stats.dislikes, 0) AS dislikes,
			COALESCE(comment_stats.comment_count, 0) AS comments,
			posts.created_at
		FROM posts 
		INNER JOIN users ON posts.id_user = users.id

		-- Sous-requête pour les likes/dislikes
		LEFT JOIN (
			SELECT 
				id_post,
				SUM(CASE WHEN liked = 1 THEN 1 ELSE 0 END) AS likes,
				SUM(CASE WHEN liked = 0 THEN 1 ELSE 0 END) AS dislikes
			FROM liked_posts
			GROUP BY id_post
		) AS like_stats ON like_stats.id_post = posts.id

		-- Sous-requête pour les commentaires
		LEFT JOIN (
			SELECT 
				id_post,
				COUNT(*) AS comment_count
			FROM comments
			GROUP BY id_post
		) AS comment_stats ON comment_stats.id_post = posts.id`

	if search == "" {
		switch order { 
		case "oldest":
			Request = BaseRequest + " ORDER BY posts.created_at ASC;"
		case "popular":
			Request = BaseRequest + " ORDER BY like_stats.likes DESC;"
		case "noresponse":
			Request = BaseRequest + " WHERE COALESCE(comment_stats.comment_count, 0) = 0;"
		default:
			Request = BaseRequest + " ORDER BY posts.created_at DESC;"
		}
	} else {
		Request = BaseRequest + " WHERE posts.title LIKE '%" + search + "%' ORDER BY posts.created_at ASC;"
	}
	

	rows, err := r.db.Query(Request)
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

func (r *PostRepository) GetPostByID(id int) (domain.Post, error) {
	var errorPost domain.Post

	rows, err := r.db.Query(`
		SELECT 
			posts.id, 
			posts.title, 
			posts.content, 
			posts.image,
			users.username,
			COALESCE(like_stats.likes, 0) AS likes,
			COALESCE(like_stats.dislikes, 0) AS dislikes,
			COALESCE(comment_stats.comment_count, 0) AS comments,
			posts.created_at
		FROM posts 
		INNER JOIN users ON posts.id_user = users.id

		-- Sous-requête pour les likes/dislikes
		LEFT JOIN (
			SELECT 
				id_post,
				SUM(CASE WHEN liked = 1 THEN 1 ELSE 0 END) AS likes,
				SUM(CASE WHEN liked = 0 THEN 1 ELSE 0 END) AS dislikes
			FROM liked_posts
			GROUP BY id_post
		) AS like_stats ON like_stats.id_post = posts.id

		-- Sous-requête pour les commentaires
		LEFT JOIN (
			SELECT 
				id_post,
				COUNT(*) AS comment_count
			FROM comments
			GROUP BY id_post
		) AS comment_stats ON comment_stats.id_post = posts.id WHERE posts.id = ?;
	`, id)
	if err != nil {
		return errorPost, err
	}
	defer rows.Close()

	var post domain.Post
	for rows.Next() {
		var createdAt string
		var imageBytes []byte
		if err := rows.Scan(&post.Id, &post.Title, &post.Content, &imageBytes, &post.Username,
			&post.Likes, &post.Dislikes, &post.Comments, &createdAt); err != nil {
			return errorPost, err
		}
		post.CreatedAt = createdAt

		if len(imageBytes) > 0 {
			post.Image = template.URL("data:image/png;base64," + base64.StdEncoding.EncodeToString(imageBytes))
		} else {
			post.Image = ""
		}

		// Récupérer les catégories pour ce post
		catRows, err := r.db.Query(`
			SELECT categories.name 
			FROM categories 
			INNER JOIN affectation ON affectation.id_category = categories.id 
			WHERE affectation.id_post = ?
		`, post.Id)
		if err != nil {
			return errorPost, err
		}

		for catRows.Next() {
			var category string
			if err := catRows.Scan(&category); err != nil {
				catRows.Close()
				return errorPost, err
			}
			post.Categories = append(post.Categories, category)
		}
		catRows.Close()
	}

	return post, nil
}

func (r *PostRepository) GetCommentsByPostID(id int) ([]domain.Comment, error) {
	rows, err := r.db.Query(`
		SELECT 
			comments.id,
			comments.content,
			users.username,
			SUM(CASE WHEN liked_comments.liked = 1 THEN 1 ELSE 0 END) AS likes,
			SUM(CASE WHEN liked_comments.liked = 0 THEN 1 ELSE 0 END) AS dislikes,
			comments.created_at
		FROM comments
		INNER JOIN users ON comments.id_user = users.id
		LEFT JOIN liked_comments ON liked_comments.id_comment = comments.id
		WHERE comments.id_post = ?
		GROUP BY comments.id;
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var comment domain.Comment
		var createdAt string
		if err := rows.Scan(&comment.Id, &comment.Content, &comment.Username,
			&comment.Likes, &comment.Dislikes, &createdAt); err != nil {
			return nil, err
		}
		comment.CreatedAt = createdAt

		comments = append(comments, comment)
	}

	return comments, nil
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
