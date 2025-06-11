package handlers

import (
	"LeForum/internal/auth/session"
	"LeForum/internal/service"
	"net/http"
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

	postID := 1

	err = h.commentService.CreateComment(
		r.FormValue("commentContent"),
		session.ID,
		postID,
	)

	if err != nil {
		http.Error(w, "Erreur lors de la création du post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}