package handlers

import (
	"LeForum/internal/auth/session"
	"LeForum/internal/service"
	"log"
	"net/http"
)

type HomeHandler struct {
	postService     *service.PostService
	categoryService *service.CategoryService
	sessionService  *session.Service
	templateService *TemplateService
}

func NewHomeHandler(ps *service.PostService, cs *service.CategoryService, ss *session.Service, ts *TemplateService) *HomeHandler {
	return &HomeHandler{
		postService:     ps,
		categoryService: cs,
		sessionService:  ss,
		templateService: ts,
	}
}

func (h *HomeHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.HomePageHandler)
	mux.HandleFunc("/theme", h.ToggleThemeHandler)
}

func getDarkModeFromCookie(r *http.Request) bool {
	cookie, err := r.Cookie("darkMode")
	if err == nil && cookie.Value == "true" {
		return true
	}
	return false
}

func (h *HomeHandler) HomePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Get current user if logged in
	user, err := h.sessionService.GetCurrentUser(r)
	if err != nil {
		log.Printf("Error getting current user: %v", err)
	}

	// Get dark mode preference
	darkMode := getDarkModeFromCookie(r)

	// Récupération des catégories
	categories, err := h.categoryService.GetCategories()
	if err != nil {
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}

	// Récupération des posts
	posts, err := h.postService.GetPosts()
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"DarkMode":      darkMode,
		"CurrentPage":   "home",
		"User":          user,
		"AllCategories": categories,
		"Posts":         posts,
	}

	err = h.templateService.RenderTemplate(w, "home_page.html", data)
	if err != nil {
		http.Error(w, "Erreur d'affichage du template: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *HomeHandler) ToggleThemeHandler(w http.ResponseWriter, r *http.Request) {
	darkMode := getDarkModeFromCookie(r)
	darkMode = !darkMode

	cookieValue := "false"
	if darkMode {
		cookieValue = "true"
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "darkMode",
		Value:    cookieValue,
		Path:     "/",
		Domain:   "forum.ynov.zeteox.fr",
		MaxAge:   365 * 24 * 60 * 60, // 1 an
		HttpOnly: false,              // Accessible par JavaScript
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	// Rediriger vers la page d'origine
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}
