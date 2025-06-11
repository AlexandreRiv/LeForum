package storage

import "LeForum/internal/domain"

func GetCategories() ([]string, error) {
	rows, err := DB.Query("SELECT name FROM categories;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func GetCategoryPosts(categoryName string) ([]domain.Post, error) {
	postRows, err := DB.Query(`
        SELECT p.id, p.title, p.content, u.username,
               SUM(CASE WHEN lp.liked = 1 THEN 1 ELSE 0 END) AS likes,
               SUM(CASE WHEN lp.liked = 0 THEN 1 ELSE 0 END) AS dislikes,
               COUNT(DISTINCT c.id) AS comments,
               p.created_at
        FROM posts p
        INNER JOIN users u ON p.id_user = u.id
        INNER JOIN affectation a ON a.id_post = p.id
        INNER JOIN categories cat ON a.id_category = cat.id
        LEFT JOIN liked_posts lp ON lp.id_post = p.id
        LEFT JOIN comments c ON c.id_post = p.id
        WHERE cat.name = ?
        GROUP BY p.id
        ORDER BY p.created_at DESC
        LIMIT 2
    `, categoryName)

	if err != nil {
		return nil, err
	}
	defer postRows.Close()

	var posts []domain.Post
	for postRows.Next() {
		var post domain.Post
		if err := postRows.Scan(&post.Id, &post.Title, &post.Content, &post.Username,
			&post.Likes, &post.Dislikes, &post.Comments, &post.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
