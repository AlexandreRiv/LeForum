package api

import (
    "LeForum/internal/auth"
    "LeForum/internal/storage"
    "html/template"
    "net/http"
    "time"
)

type PageData struct {
    DarkMode      bool
    CurrentPage   string
    Posts         []storage.Post
    AllCategories []string
    User          *auth.LoggedUser
}

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
    // Get current user if logged in
    user, _ := auth.GetCurrentUser(r)
    
    data := PageData{
        DarkMode:    getDarkModeFromCookie(r),
        CurrentPage: "categories",
        User:        user,
    }
    
    tmpl := template.Must(template.ParseFiles("web/templates/categories.html"))
    template.Must(tmpl.ParseGlob("web/templates/components/*.html"))
    tmpl.Execute(w, data)
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