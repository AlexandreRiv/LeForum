package storage

import (
	"LeForum/internal/domain"
	"time"
)

func CreatePost(title, content, sessionID, category string) error {
	createdAt := time.Now().Add(2 * time.Hour)

	_, err := DB.Exec(
		"INSERT INTO posts (title, content, id_user, created_at) VALUES (?, ?, (SELECT users.id FROM users INNER JOIN sessions ON users.mail = sessions.user_email WHERE sessions.id = ?), ?);",
		title,
		content,
		sessionID,
		createdAt,
	)
	if err != nil {
		return err
	}

	_, err = DB.Exec(
		"INSERT INTO affectation VALUES ((SELECT id FROM posts WHERE title = ? AND content = ? AND id_user = (SELECT users.id FROM users INNER JOIN sessions ON users.mail = sessions.user_email WHERE sessions.id = ?)), (SELECT id FROM categories WHERE name = ?));",
		title,
		content,
		sessionID,
		category,
	)

	return err
}

func GetPosts() ([]domain.Post, error) {
	var posts []domain.Post

	rows, err := DB.Query("SELECT posts.id,posts.title,posts.content,users.username,SUM(CASE WHEN liked_posts.liked = 1 THEN 1 ELSE 0 END) AS likes,SUM(CASE WHEN liked_posts.liked = 0 THEN 1 ELSE 0 END) AS dislikes,COUNT(distinct comments.id) AS comments FROM posts INNER JOIN users ON posts.id_user = users.id LEFT JOIN liked_posts ON liked_posts.id_post = posts.id LEFT JOIN comments ON comments.id_post = posts.id GROUP BY posts.id;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post domain.Post
		if err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.Username, &post.Likes, &post.Dislikes, &post.Comments); err != nil {
			return nil, err
		}

		catRows, err := DB.Query("SELECT categories.name FROM categories INNER JOIN affectation ON affectation.id_category = categories.id INNER JOIN posts ON posts.id = affectation.id_post WHERE posts.id = ?;", post.Id)
		if err != nil {
			return nil, err
		}
		defer catRows.Close()

		for catRows.Next() {
			var category string
			if err := catRows.Scan(&category); err != nil {
				return nil, err
			}
			post.Categories = append(post.Categories, category)
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func LikePost(sessionID string, postID string, likeType int) error {
	// Tente d'abord d'insérer le like
	_, err := DB.Exec(
		"INSERT INTO liked_posts VALUES ((SELECT users.id FROM users INNER JOIN sessions ON users.mail = sessions.user_email WHERE sessions.id = ?), ?, ?);",
		sessionID,
		postID,
		likeType,
	)

	// Si l'insertion échoue, essaye de supprimer
	if err != nil {
		_, err = DB.Exec(
			"DELETE FROM liked_posts WHERE id_user = (SELECT users.id FROM users INNER JOIN sessions ON users.mail = sessions.user_email WHERE sessions.id = ?) AND id_post = ?;",
			sessionID,
			postID,
		)
	}

	return err
}
