package api

import (
	"LeForum/internal/auth"
	"LeForum/internal/storage"
	"html/template"
	"net/http"
	"strings"
	"time"
)

type PageData struct {
	DarkMode      bool
	CurrentPage   string
	Posts         []storage.Post
	AllCategories []string
	User          *auth.LoggedUser
}

var templateFuncs = template.FuncMap{
	"ToLower": strings.ToLower,
}

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	// Get current user if logged in
	user, _ := auth.GetCurrentUser(r)

	// Définir les fonctions template
	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
		"Mod":     func(i, j int) int { return i % j },
	}

	// Données de la page
	data := PageData{
		DarkMode:    getDarkModeFromCookie(r),
		CurrentPage: "categories",
		User:        user,
	}

	// Définir des styles prédéfinis pour une rotation
	gradientStyles := []struct {
		Icon         string
		GradientFrom string
		GradientTo   string
	}{
		{"fi fi-rr-star", "from-blue-500", "to-indigo-500"},
		{"fi fi-rr-comment", "from-green-500", "to-teal-500"},
		{"fi fi-rr-book", "from-red-500", "to-orange-500"},
		{"fi fi-rr-world", "from-purple-500", "to-pink-500"},
		{"fi fi-rr-puzzle-piece", "from-yellow-500", "to-amber-500"},
		{"fi fi-rr-diamond", "from-cyan-500", "to-blue-500"},
	}

	// Récupérer toutes les catégories
	rows, err := storage.DB.Query("SELECT name FROM categories")
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des catégories", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Liste pour stocker les catégories avec leurs posts
	var categories []struct {
		Name         string
		Icon         string
		GradientFrom string
		GradientTo   string
		Posts        []storage.Post
	}

	// Pour chaque catégorie, récupérer les 2 derniers posts
	i := 0
	for rows.Next() {
		var categoryName string
		if err := rows.Scan(&categoryName); err != nil {
			http.Error(w, "Erreur lors de la lecture des catégories", http.StatusInternalServerError)
			return
		}

		// Sélectionner un style basé sur l'index (rotation cyclique)
		styleIndex := i % len(gradientStyles)
		style := gradientStyles[styleIndex]

		// Requête pour récupérer les 2 derniers posts de cette catégorie
		postRows, err := storage.DB.Query(`
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
			http.Error(w, "Erreur lors de la récupération des posts", http.StatusInternalServerError)
			return
		}

		var posts []storage.Post
		for postRows.Next() {
			var post storage.Post
			if err := postRows.Scan(&post.Id, &post.Title, &post.Content, &post.Username,
				&post.Likes, &post.Dislikes, &post.Comments, &post.CreatedAt); err != nil {
				postRows.Close()
				http.Error(w, "Erreur lors de la lecture des posts", http.StatusInternalServerError)
				return
			}
			posts = append(posts, post)
		}
		postRows.Close()

		// Ajouter cette catégorie avec ses posts
		categories = append(categories, struct {
			Name         string
			Icon         string
			GradientFrom string
			GradientTo   string
			Posts        []storage.Post
		}{
			Name:         categoryName,
			Icon:         style.Icon,
			GradientFrom: style.GradientFrom,
			GradientTo:   style.GradientTo,
			Posts:        posts,
		})

		i++
	}

	// Ajouter les catégories au data
	data.AllCategories = make([]string, len(categories))
	for i, cat := range categories {
		data.AllCategories[i] = cat.Name
	}

	// Créer le template avec son nom exact et les fonctions avant le parsing
	tmpl := template.New("categories.html").Funcs(funcMap)

	// Parser le fichier principal
	tmpl, err = tmpl.ParseFiles("web/templates/categories.html")
	if err != nil {
		http.Error(w, "Erreur de chargement du template principal: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Parser les composants
	tmpl, err = tmpl.ParseGlob("web/templates/components/*.html")
	if err != nil {
		http.Error(w, "Erreur de chargement des templates de composants: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Exécuter le template
	if err := tmpl.Execute(w, struct {
		PageData
		Categories []struct {
			Name         string
			Icon         string
			GradientFrom string
			GradientTo   string
			Posts        []storage.Post
		}
	}{
		PageData:   data,
		Categories: categories,
	}); err != nil {
		http.Error(w, "Erreur d'affichage du template: "+err.Error(), http.StatusInternalServerError)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// Get current user if logged in
	user, _ := auth.GetCurrentUser(r)

	// check if the dark mode cookie exists
	darkMode := getDarkModeFromCookie(r)

	data := PageData{
		DarkMode:    darkMode,
		CurrentPage: "home",
		User:        user,
	}

	tmpl := template.Must(template.ParseFiles("web/templates/home_page.html"))
	template.Must(tmpl.ParseGlob("web/templates/components/*.html"))

	SQLCatRequest := "SELECT name FROM categories;"
	SQLPostsRequest := "SELECT posts.id,posts.title,posts.content,users.username,SUM(CASE WHEN liked_posts.liked = 1 THEN 1 ELSE 0 END) AS likes,SUM(CASE WHEN liked_posts.liked = 0 THEN 1 ELSE 0 END) AS dislikes,COUNT(distinct comments.id) AS comments FROM posts INNER JOIN users ON posts.id_user = users.id LEFT JOIN liked_posts ON liked_posts.id_post = posts.id LEFT JOIN comments ON comments.id_post = posts.id GROUP BY posts.id;"
	SQLPostsCatRequest := "SELECT categories.name FROM categories INNER JOIN affectation ON affectation.id_category = categories.id INNER JOIN posts ON posts.id = affectation.id_post WHERE posts.id = ?;"

	rows, err := storage.DB.Query(SQLCatRequest)
	if err != nil {
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			http.Error(w, "Failed to scan category", http.StatusInternalServerError)
			return
		}
		data.AllCategories = append(data.AllCategories, category)
	}

	rows, err = storage.DB.Query(SQLPostsRequest)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post storage.Post
		if err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.Username, &post.Likes, &post.Dislikes, &post.Comments); err != nil {
			http.Error(w, "Failed to scan post", http.StatusInternalServerError)
			return
		}

		catRows, err := storage.DB.Query(SQLPostsCatRequest, post.Id)
		if err != nil {
			http.Error(w, "Failed to fetch post categories", http.StatusInternalServerError)
			return
		}
		defer catRows.Close()

		for catRows.Next() {
			var category string
			if err := catRows.Scan(&category); err != nil {
				http.Error(w, "Failed to scan post category", http.StatusInternalServerError)
				return
			}
			post.Categories = append(post.Categories, category)
		}

		data.Posts = append(data.Posts, post)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Erreur d'affichage du template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func ToggleThemeHandler(w http.ResponseWriter, r *http.Request) {
	// read the cookie to check the current theme
	darkMode := getDarkModeFromCookie(r)

	// reverse the theme
	darkMode = !darkMode

	// set the cookie to 1 year for the expiration
	expiration := time.Now().Add(365 * 24 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:    "darkMode",
		Value:   boolToString(darkMode),
		Expires: expiration,
		Path:    "/",
	})

	// redirect to the referer page
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

func getDarkModeFromCookie(r *http.Request) bool {
	cookie, err := r.Cookie("darkMode")
	if err == nil && cookie.Value == "true" {
		return true
	}
	return false
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
