package handlers

import (
	"LeForum/internal/auth"
	"LeForum/internal/domain"
	"LeForum/internal/storage"
	"html/template"
	"net/http"
	"strings"
)

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	// Get current user if logged in
	user, _ := auth.GetCurrentUser(r)

	// Récupérer toutes les catégories
	categoriesNames, err := storage.GetCategories()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des catégories", http.StatusInternalServerError)
		return
	}

	// Données de la page
	data := struct {
		DarkMode      bool
		CurrentPage   string
		User          *auth.LoggedUser
		AllCategories []string
	}{
		DarkMode:      getDarkModeFromCookie(r),
		CurrentPage:   "categories",
		User:          user,
		AllCategories: categoriesNames,
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

	// Liste pour stocker les catégories avec leurs posts
	var categories []struct {
		Name         string
		NameLower    string
		Icon         string
		GradientFrom string
		GradientTo   string
		Posts        []domain.Post
	}

	// Pour chaque catégorie, récupérer les 2 derniers posts
	for i, categoryName := range categoriesNames {
		// Sélectionner un style basé sur l'index (rotation cyclique)
		styleIndex := i % len(gradientStyles)
		style := gradientStyles[styleIndex]

		// Récupérer les posts pour cette catégorie
		posts, err := storage.GetCategoryPosts(categoryName)
		if err != nil {
			http.Error(w, "Erreur lors de la récupération des posts", http.StatusInternalServerError)
			return
		}

		// Ajouter cette catégorie avec ses posts
		categories = append(categories, struct {
			Name         string
			NameLower    string
			Icon         string
			GradientFrom string
			GradientTo   string
			Posts        []domain.Post
		}{
			Name:         categoryName,
			NameLower:    strings.ToLower(categoryName),
			Icon:         style.Icon,
			GradientFrom: style.GradientFrom,
			GradientTo:   style.GradientTo,
			Posts:        posts,
		})
	}

	// Parser les templates
	tmpl := template.New("categories.html")
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
		DarkMode      bool
		CurrentPage   string
		User          *auth.LoggedUser
		AllCategories []string
		Categories    []struct {
			Name         string
			NameLower    string
			Icon         string
			GradientFrom string
			GradientTo   string
			Posts        []domain.Post
		}
	}{
		DarkMode:      data.DarkMode,
		CurrentPage:   data.CurrentPage,
		User:          data.User,
		AllCategories: data.AllCategories,
		Categories:    categories,
	}); err != nil {
		http.Error(w, "Erreur d'affichage du template: "+err.Error(), http.StatusInternalServerError)
	}
}
