package handlers

import (
	"LeForum/internal/auth/session"
	"LeForum/internal/service"
	"io"
	"log"
	"net/http"
	"strconv"
)

type PostHandler struct {
	postService     *service.PostService
	categoryService *service.CategoryService
	sessionService  *session.Service
	templateService *TemplateService
}

func NewPostHandler(ps *service.PostService, cs *service.CategoryService, ss *session.Service, ts *TemplateService) *PostHandler {
	return &PostHandler{
		postService:     ps,
		categoryService: cs,
		sessionService:  ss,
		templateService: ts,
	}
}

func (h *PostHandler) PostPageHandler(w http.ResponseWriter, r *http.Request) {
	user, err := h.sessionService.GetCurrentUser(r)
	if err != nil {
		log.Printf("Error getting user: %v", err)
	}

	PotsIDStr := r.URL.Query().Get("id")
	PostID, err := strconv.Atoi(PotsIDStr)
	if err != nil {
		http.Error(w, "Like parameter is invalid", http.StatusBadRequest)
		return
	}

	post, _ := h.postService.GetPostByID(PostID)

	comments, err := h.postService.GetCommentsByPostID(PostID)
	if err != nil {
		log.Printf("Erreur lors de la récupération des commentaires du post %d: %v", PostID, err)
		comments = nil
	}

	// Récupération des catégories
	categories, err := h.categoryService.GetCategories()
	if err != nil {
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}

	darkMode := getDarkModeFromCookie(r)

	data := map[string]interface{}{
		"DarkMode":      darkMode,
		"CurrentPage":   "post",
		"AllCategories": categories,
		"User":          user,
		"Post":          post,
		"Comments":      comments,
	}

	err = h.templateService.RenderTemplate(w, "post_page.html", data)
	if err != nil {
		http.Error(w, "Erreur d'affichage du template: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *PostHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, err := h.sessionService.GetSession(r)
	if err != nil || session == nil {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	var imageBytes []byte
	file, _, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		imageBytes, err = io.ReadAll(file)
		if err != nil {
			http.Error(w, "Erreur lors de la lecture des bytes de l'image", http.StatusInternalServerError)
			return
		}
	} else if err != http.ErrMissingFile {
		http.Error(w, "Erreur lors de la lecture du fichier image", http.StatusInternalServerError)
		return
	}

	err = h.postService.CreatePost(
		r.FormValue("title"),
		r.FormValue("content"),
		session.ID,
		r.FormValue("category"),
		imageBytes,
	)

	if err != nil {
		http.Error(w, "Erreur lors de la création du post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *PostHandler) LikePostHandler(w http.ResponseWriter, r *http.Request) {
	session, err := h.sessionService.GetSession(r)
	if err != nil || session == nil {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	postID := r.URL.Query().Get("id")
	if postID == "" {
		http.Error(w, "Id parameter is missing", http.StatusBadRequest)
		return
	}

	likeTypeStr := r.URL.Query().Get("like")
	likeType, err := strconv.Atoi(likeTypeStr)
	if err != nil {
		http.Error(w, "Like parameter is invalid", http.StatusBadRequest)
		return
	}

	err = h.postService.LikePost(session.ID, postID, likeType)
	if err != nil {
		http.Error(w, "Like error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *PostHandler) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	PotsIDStr := r.URL.Query().Get("id")
	PostID, err := strconv.Atoi(PotsIDStr)
	if err != nil {
		http.Error(w, "Like parameter is invalid", http.StatusBadRequest)
		return
	}

	session, err := h.sessionService.GetSession(r)
	if err != nil || session == nil {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	err = h.postService.UpdatePost(
		PostID,
		r.FormValue("title"),
		r.FormValue("content"),
		r.FormValue("category"),
	)

	if err != nil {
		http.Error(w, "Erreur lors de la création du post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post?id="+PotsIDStr, http.StatusSeeOther)
}

func (h *PostHandler) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	PotsIDStr := r.URL.Query().Get("id")
	PostID, err := strconv.Atoi(PotsIDStr)
	if err != nil {
		http.Error(w, "Like parameter is invalid", http.StatusBadRequest)
		return
	}

	err = h.postService.DeletePost(PostID)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *PostHandler) EditPostHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID du post à éditer
	postIDStr := r.URL.Query().Get("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "ID de post invalide", http.StatusBadRequest)
		return
	}

	// Récupérer l'utilisateur actuel
	user, err := h.sessionService.GetCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	// Récupérer les informations du post
	post, err := h.postService.GetPostByID(postID)
	if err != nil {
		http.Error(w, "Impossible de récupérer les informations du post", http.StatusInternalServerError)
		return
	}

	// Récupérer les catégories disponibles
	categories, err := h.categoryService.GetCategories()
	if err != nil {
		http.Error(w, "Impossible de récupérer les catégories", http.StatusInternalServerError)
		return
	}

	darkMode := getDarkModeFromCookie(r)

	data := map[string]interface{}{
		"User":          user,
		"Post":          post,
		"AllCategories": categories,
		"DarkMode":      darkMode,
		"PageTitle":     "Modifier le post",
	}

	err = h.templateService.RenderTemplate(w, "edit_post.html", data)
	if err != nil {
		http.Error(w, "Erreur de rendu du template: "+err.Error(), http.StatusInternalServerError)
	}
}
