package handlers

import (
	"LeForum/internal/auth/session"
	"LeForum/internal/service"
	"log"
	"net/http"
	"strconv"
	"io"
)

type PostHandler struct {
	postService     *service.PostService
	sessionService  *session.Service
	templateService *TemplateService
}

func NewPostHandler(ps *service.PostService, ss *session.Service, ts *TemplateService) *PostHandler {
	return &PostHandler{
		postService:     ps,
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

	post, err := h.postService.GetPostByID(PostID)
	if err != nil {
		log.Printf("Erreur lors de la récupération du post %d: %v", PostID, err)
		http.Error(w, "Post introuvable", http.StatusNotFound)
		return
	}

	comments, err := h.postService.GetCommentsByPostID(PostID)
	if err != nil {
		log.Printf("Erreur lors de la récupération des commentaires du post %d: %v", PostID, err)
		comments = nil 
	}

	darkMode := getDarkModeFromCookie(r)

	data := map[string]interface{}{
		"DarkMode":    darkMode,
		"CurrentPage": "post",
		"User":        user,
		"Post":	   	   post,
		"Comments":	   comments,
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

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Erreur lors de la lecture du fichier image",  http.StatusMethodNotAllowed)
	}
	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Erreur lors de la lecture des bytes de l'image",  http.StatusMethodNotAllowed)
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
