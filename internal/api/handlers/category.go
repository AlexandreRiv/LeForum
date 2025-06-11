package handlers

import (
	"LeForum/internal/auth/session"
	"LeForum/internal/domain"
	"LeForum/internal/service"
	"log"
	"net/http"
	"strings"
)

type CategoryHandler struct {
	categoryService *service.CategoryService
	sessionService  *session.Service
	templateService *TemplateService
}

func NewCategoryHandler(cs *service.CategoryService, ss *session.Service, ts *TemplateService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: cs,
		sessionService:  ss,
		templateService: ts,
	}
}

func (h *CategoryHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/categories", h.CategoriesHandler)
}

func (h *CategoryHandler) CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	// Get current user if logged in
	user, err := h.sessionService.GetCurrentUser(r)
	if err != nil {
		log.Printf("Erreur récupération utilisateur: %v", err)
	}

	// Récupérer toutes les catégories
	categoriesNames, err := h.categoryService.GetCategories()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des catégories", http.StatusInternalServerError)
		return
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
		Posts        []interface{}
	}

	// Pour chaque catégorie, récupérer les 2 derniers posts
	for i, categoryName := range categoriesNames {
		// Sélectionner un style basé sur l'index (rotation cyclique)
		styleIndex := i % len(gradientStyles)
		style := gradientStyles[styleIndex]

		// Récupérer les posts pour cette catégorie
		posts, err := h.categoryService.GetCategoryPosts(categoryName)
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
			Posts        []interface{}
		}{
			Name:         categoryName,
			NameLower:    strings.ToLower(categoryName),
			Icon:         style.Icon,
			GradientFrom: style.GradientFrom,
			GradientTo:   style.GradientTo,
			Posts:        interfaceSlice(posts),
		})
	}

	// Données de la page
	data := map[string]interface{}{
		"DarkMode":      getDarkModeFromCookie(r),
		"CurrentPage":   "categories",
		"User":          user,
		"AllCategories": categoriesNames,
		"Categories":    categories,
	}

	err = h.templateService.RenderTemplate(w, "categories.html", data)
	if err != nil {
		http.Error(w, "Erreur d'affichage du template: "+err.Error(), http.StatusInternalServerError)
	}
}

// Helper function to convert []domain.Post to []interface{}
func interfaceSlice(posts interface{}) []interface{} {
	switch posts := posts.(type) {
	case []domain.Post:
		result := make([]interface{}, len(posts))
		for i, v := range posts {
			result[i] = v
		}
		return result
	default:
		return []interface{}{}
	}
}
