package handlers

import (
	"LeForum/internal/domain"
	"LeForum/internal/auth/session"
	"LeForum/internal/service"
	"log"
	"net/http"
)

type HomeHandler struct {
	postService     *service.PostService
	categoryService *service.CategoryService
	notificationService *service.NotificationService
	sessionService  *session.Service
	templateService *TemplateService
}

func NewHomeHandler(ps *service.PostService, cs *service.CategoryService, ns *service.NotificationService, ss *session.Service, ts *TemplateService) *HomeHandler {
	return &HomeHandler{
		postService:     ps,
		categoryService: cs,
		notificationService: ns,
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
	order := r.URL.Query().Get("filter")
	search := r.URL.Query().Get("search")

	if order == "" {
		order = "newest" // Default order
	}

	posts, err := h.postService.GetPosts(order, search)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}


	var notifs []domain.Notification
	session, err := h.sessionService.GetSession(r)
	if err != nil {
		http.Error(w, "Error fetching sessions", http.StatusInternalServerError)
		return
	}
	if session != nil {
		notifs, err = h.notificationService.GetNotifications(session.ID)
	}

	data := map[string]interface{}{
		"DarkMode":      darkMode,
		"CurrentPage":   "home",
		"User":          user,
		"AllCategories": categories,
		"Posts":         posts,
		"ActiveFilter":  order,
		"Notifications": notifs,
		"NotificationNb": len(notifs),
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
