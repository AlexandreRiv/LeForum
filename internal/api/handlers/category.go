package handlers

import (
	"LeForum/internal/auth/session"
	"LeForum/internal/domain"
	"LeForum/internal/service"
	"log"
	"net/http"
	"strconv"
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
	mux.HandleFunc("/category/", h.CategoryPostsHandler)
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

	// Liste pour stocker les catégories avec leurs posts
	var categories []struct {
		Name      string
		NameLower string
		Posts     []interface{}
	}

	// Pour chaque catégorie, récupérer les 2 derniers posts
	for _, categoryName := range categoriesNames {
		// Récupérer les posts pour cette catégorie
		posts, err := h.categoryService.GetCategoryPosts(categoryName)
		if err != nil {
			http.Error(w, "Erreur lors de la récupération des posts", http.StatusInternalServerError)
			return
		}

		// Ajouter cette catégorie avec ses posts
		categories = append(categories, struct {
			Name      string
			NameLower string
			Posts     []interface{}
		}{
			Name:      categoryName,
			NameLower: strings.ToLower(categoryName),
			Posts:     interfaceSlice(posts),
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

func (h *CategoryHandler) CategoryPostsHandler(w http.ResponseWriter, r *http.Request) {
	// Extraire le nom de la catégorie de l'URL (/category/nom-categorie)
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		http.Error(w, "Catégorie non spécifiée", http.StatusBadRequest)
		return
	}

	categoryNameLower := parts[2]
	// Convertir le premier caractère en majuscule pour la recherche dans la base de données
	categoryName := strings.Title(categoryNameLower)

	// Récupérer l'utilisateur actuel
	user, err := h.sessionService.GetCurrentUser(r)
	if err != nil {
		log.Printf("Erreur récupération utilisateur: %v", err)
	}

	// Gérer la pagination
	page := 1
	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		pageNum, err := strconv.Atoi(pageStr)
		if err == nil && pageNum > 0 {
			page = pageNum
		}
	}
	postsPerPage := 10
	offset := (page - 1) * postsPerPage

	// Récupérer tous les posts de cette catégorie avec pagination
	posts, total, err := h.categoryService.GetAllCategoryPosts(categoryName, postsPerPage, offset)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des posts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Récupérer toutes les catégories pour le menu
	allCategories, err := h.categoryService.GetCategories()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des catégories", http.StatusInternalServerError)
		return
	}

	// Calculer les infos de pagination
	totalPages := (total + postsPerPage - 1) / postsPerPage // Arrondir au supérieur
	showPagination := totalPages > 1

	// Générer les liens de pagination
	var paginationLinks []struct {
		Page   int
		Active bool
	}

	maxVisiblePages := 5
	startPage := page - maxVisiblePages/2
	if startPage < 1 {
		startPage = 1
	}

	endPage := startPage + maxVisiblePages - 1
	if endPage > totalPages {
		endPage = totalPages
		startPage = endPage - maxVisiblePages + 1
		if startPage < 1 {
			startPage = 1
		}
	}

	for i := startPage; i <= endPage; i++ {
		paginationLinks = append(paginationLinks, struct {
			Page   int
			Active bool
		}{
			Page:   i,
			Active: i == page,
		})
	}

	// Données de la page
	data := map[string]interface{}{
		"DarkMode":          getDarkModeFromCookie(r),
		"CurrentPage":       "categories",
		"User":              user,
		"AllCategories":     allCategories,
		"CategoryName":      categoryName,
		"CategoryNameLower": categoryNameLower,
		"Posts":             interfaceSlice(posts),
		"Page":              page,
		"TotalPages":        totalPages,
		"PrevPage":          max(1, page-1),
		"NextPage":          min(totalPages, page+1),
		"ShowPagination":    showPagination,
		"PaginationLinks":   paginationLinks,
	}

	err = h.templateService.RenderTemplate(w, "category_posts.html", data)
	if err != nil {
		http.Error(w, "Erreur d'affichage du template: "+err.Error(), http.StatusInternalServerError)
	}
}

// Fonctions utilitaires
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
