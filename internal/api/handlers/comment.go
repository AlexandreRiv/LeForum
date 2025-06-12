package handlers

import (
	"LeForum/internal/auth/session"
	"LeForum/internal/service"
	"net/http"
	"strconv"
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

	PotsIDStr := r.URL.Query().Get("id")
	PostID, err := strconv.Atoi(PotsIDStr)
	if err != nil {
		http.Error(w, "Like parameter is invalid", http.StatusBadRequest)
		return
	}

	err = h.commentService.CreateComment(
		r.FormValue("commentContent"),
		session.ID,
		PostID,
	)

	if err != nil {
		http.Error(w, "Erreur lors de la création du post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post?id="+PotsIDStr, http.StatusSeeOther)
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