package handlers

import (
	"LeForum/internal/auth/session"
	"LeForum/internal/service"
	"net/http"
	"strconv"
	"io"
	"fmt"
)

type CommentHandler struct {
	commentService  *service.CommentService
	sessionService  *session.Service
	templateService *TemplateService
}

func NewCommentHandler(cs *service.CommentService, ss *session.Service, ts *TemplateService) *CommentHandler {
	return &CommentHandler{
		commentService:  cs,
		sessionService:  ss,
		templateService: ts,
	}
}

func (h *CommentHandler) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, err := h.sessionService.GetSession(r)
	if err != nil || session == nil {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	PostIDStr := r.URL.Query().Get("id")
	PostID, err := strconv.Atoi(PostIDStr)
	if err != nil {
		http.Error(w, "Like parameter is invalid", http.StatusBadRequest)
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

	err = h.commentService.CreateComment(
		r.FormValue("commentContent"),
		session.ID,
		PostID,
		imageBytes,
	)

	if err != nil {
		fmt.Printf("Erreur lors de la création du commentaire : %v\n", err)
		http.Error(w, "Erreur lors de la création du post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post?id="+PostIDStr, http.StatusSeeOther)
}

func (h *CommentHandler) LikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	session, err := h.sessionService.GetSession(r)
	if err != nil || session == nil {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	commentID := r.URL.Query().Get("id")
	if commentID == "" {
		http.Error(w, "Id parameter is missing", http.StatusBadRequest)
		return
	}

	likeTypeStr := r.URL.Query().Get("like")
	likeType, err := strconv.Atoi(likeTypeStr)
	if err != nil {
		http.Error(w, "Like parameter is invalid", http.StatusBadRequest)
		return
	}

	PostID := r.URL.Query().Get("postId")

	err = h.commentService.LikeComment(session.ID, commentID, likeType)
	if err != nil {
		http.Error(w, "Like error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post?id="+PostID, http.StatusSeeOther)
}

func (h *CommentHandler) DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	CommentIDStr := r.URL.Query().Get("id")
	CommentID, err := strconv.Atoi(CommentIDStr)
	if err != nil {
		http.Error(w, "Id parameter is invalid", http.StatusBadRequest)
		return
	}

	PostIDStr := r.URL.Query().Get("postId")

	err = h.commentService.DeleteComment(CommentID)

	http.Redirect(w, r, "/post?id="+PostIDStr, http.StatusSeeOther)
}