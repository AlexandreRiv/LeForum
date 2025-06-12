// internal/api/handlers/moderation.go
package handlers

import (
	"LeForum/internal/auth/session"
	"LeForum/internal/domain"
	"LeForum/internal/service"
	"log"
	"net/http"
	"strconv"
)

type ModerationHandler struct {
	reportService   *service.ReportService
	sessionService  *session.Service
	postService     *service.PostService
	commentService  *service.CommentService
	templateService *TemplateService
}

func NewModerationHandler(
	reportService *service.ReportService,
	sessionService *session.Service,
	postService *service.PostService,
	commentService *service.CommentService,
	templateService *TemplateService,
) *ModerationHandler {
	return &ModerationHandler{
		reportService:   reportService,
		sessionService:  sessionService,
		postService:     postService,
		commentService:  commentService,
		templateService: templateService,
	}
}

func (h *ModerationHandler) ModerationDashboard(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'utilisateur actuel pour afficher son nom dans la navbar
	user, _ := h.sessionService.GetCurrentUser(r)

	// Récupérer tous les signalements en attente
	reports, err := h.reportService.GetPendingReports()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des signalements", http.StatusInternalServerError)
		return
	}

	darkMode := getDarkModeFromCookie(r)

	data := map[string]interface{}{
		"User":      user,
		"Reports":   reports,
		"DarkMode":  darkMode,
		"PageTitle": "Modération",
	}

	if err := h.templateService.RenderTemplate(w, "moderation/dashboard.html", data); err != nil {
		http.Error(w, "Erreur de rendu du template: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *ModerationHandler) ReportContentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer l'utilisateur actuel
	user, err := h.sessionService.GetCurrentUser(r)
	if err != nil || user == nil {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	// Récupérer les données du formulaire
	postIDStr := r.FormValue("post_id")
	commentIDStr := r.FormValue("comment_id")
	reason := r.FormValue("reason")
	reportType := r.FormValue("type")

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "ID de post invalide", http.StatusBadRequest)
		return
	}

	// Vérifier si c'est un commentaire ou un post qui est signalé
	var commentID *int
	if commentIDStr != "" {
		id, err := strconv.Atoi(commentIDStr)
		if err != nil {
			http.Error(w, "ID de commentaire invalide", http.StatusBadRequest)
			return
		}
		commentID = &id
	}

	// Déterminer le type de signalement
	var reportTypeEnum domain.ReportType
	switch reportType {
	case "inappropriate":
		reportTypeEnum = domain.ReportContentInappropriate
	case "spam":
		reportTypeEnum = domain.ReportSpam
	case "harassment":
		reportTypeEnum = domain.ReportHarassment
	default:
		reportTypeEnum = domain.ReportOther
	}

	// Créer le signalement
	_, err = h.reportService.CreateReport(postID, commentID, user.Email, reason, reportTypeEnum)
	if err != nil {
		log.Printf("Erreur lors de la création du signalement: %v", err)
		http.Error(w, "Erreur lors de la création du signalement", http.StatusInternalServerError)
		return
	}

	// Rediriger vers la page du post
	http.Redirect(w, r, "/post?id="+postIDStr, http.StatusSeeOther)
}

func (h *ModerationHandler) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer l'ID du post à supprimer
	postIDStr := r.FormValue("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "ID de post invalide", http.StatusBadRequest)
		return
	}

	// Supprimer le post
	if err := h.postService.DeletePost(postID); err != nil {
		http.Error(w, "Erreur lors de la suppression du post", http.StatusInternalServerError)
		return
	}

	// Rediriger vers le tableau de bord de modération
	http.Redirect(w, r, "/moderation", http.StatusSeeOther)
}

func (h *ModerationHandler) DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer l'ID du commentaire à supprimer
	commentIDStr := r.FormValue("comment_id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "ID de commentaire invalide", http.StatusBadRequest)
		return
	}

	// Implémenter la suppression du commentaire dans le service
	// Note: il faut créer cette méthode dans CommentService
	if err := h.commentService.DeleteComment(commentID); err != nil {
		http.Error(w, "Erreur lors de la suppression du commentaire", http.StatusInternalServerError)
		return
	}

	// Rediriger vers le tableau de bord de modération
	http.Redirect(w, r, "/moderation", http.StatusSeeOther)
}
