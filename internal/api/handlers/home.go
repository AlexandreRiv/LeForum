package handlers

import (
	"LeForum/internal/auth"
	"LeForum/internal/domain"
	"LeForum/internal/storage"
	"html/template"
	"net/http"
)

func getDarkModeFromCookie(r *http.Request) bool {
	cookie, err := r.Cookie("darkMode")
	if err == nil && cookie.Value == "true" {
		return true
	}
	return false
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Get current user if logged in
	user, _ := auth.GetCurrentUser(r)

	// check if the dark mode cookie exists
	darkMode := getDarkModeFromCookie(r)

	// Récupération des catégories
	categories, err := storage.GetCategories()
	if err != nil {
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}

	// Récupération des posts
	posts, err := storage.GetPosts()
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	data := struct {
		DarkMode      bool
		CurrentPage   string
		User          *auth.LoggedUser
		AllCategories []string
		Posts         []domain.Post
	}{
		DarkMode:      darkMode,
		CurrentPage:   "home",
		User:          user,
		AllCategories: categories,
		Posts:         posts,
	}

	tmpl := template.Must(template.ParseFiles("web/templates/home_page.html"))
	template.Must(tmpl.ParseGlob("web/templates/components/*.html"))

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Erreur d'affichage du template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
