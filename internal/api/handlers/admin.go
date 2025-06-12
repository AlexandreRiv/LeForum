package handlers

import (
	"LeForum/internal/auth/session"
	"LeForum/internal/domain"
	"LeForum/internal/service"
	"log"
	"net/http"
	"strconv"
)

type AdminHandler struct {
	userService     *service.UserService
	categoryService *service.CategoryService
	reportService   *service.ReportService
	sessionService  *session.Service
	templateService *TemplateService
}

func NewAdminHandler(
	userService *service.UserService,
	categoryService *service.CategoryService,
	reportService *service.ReportService,
	sessionService *session.Service,
	templateService *TemplateService,
) *AdminHandler {
	return &AdminHandler{
		userService:     userService,
		categoryService: categoryService,
		reportService:   reportService,
		sessionService:  sessionService,
		templateService: templateService,
	}
}

func (h *AdminHandler) AdminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := h.sessionService.GetCurrentUser(r)

	// TODO: Ajouter des statistiques réelles ici
	data := map[string]interface{}{
		"User":          user,
		"UserCount":     0,
		"PostCount":     0,
		"CommentCount":  0,
		"CategoryCount": 0,
		"DarkMode":      getDarkModeFromCookie(r),
		"PageTitle":     "Administration",
	}

	if err := h.templateService.RenderTemplate(w, "admin/dashboard.html", data); err != nil {
		http.Error(w, "Erreur de rendu du template: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *AdminHandler) ManageUsersHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := h.sessionService.GetCurrentUser(r)

	users, err := h.userService.GetAllUsers()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des utilisateurs", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"User":      user,
		"Users":     users,
		"DarkMode":  getDarkModeFromCookie(r),
		"PageTitle": "Gestion des utilisateurs",
	}

	if err := h.templateService.RenderTemplate(w, "admin/users.html", data); err != nil {
		http.Error(w, "Erreur de rendu du template: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *AdminHandler) ChangeUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.FormValue("user_id")
	role := r.FormValue("role")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "ID utilisateur invalide", http.StatusBadRequest)
		return
	}

	// Vérifier que le rôle est valide
	var roleType domain.RoleType
	switch role {
	case "user":
		roleType = domain.RoleUser
	case "moderator":
		roleType = domain.RoleModerator
	case "admin":
		roleType = domain.RoleAdmin
	default:
		http.Error(w, "Rôle invalide", http.StatusBadRequest)
		return
	}

	if err := h.userService.UpdateUserRole(userID, roleType); err != nil {
		http.Error(w, "Erreur lors de la mise à jour du rôle", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

func (h *AdminHandler) ManageCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := h.sessionService.GetCurrentUser(r)

	categories, err := h.categoryService.GetCategories()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des catégories", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"User":       user,
		"Categories": categories,
		"DarkMode":   getDarkModeFromCookie(r),
		"PageTitle":  "Gestion des catégories",
	}

	if err := h.templateService.RenderTemplate(w, "admin/categories.html", data); err != nil {
		http.Error(w, "Erreur de rendu du template: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *AdminHandler) AddCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	categoryName := r.FormValue("category_name")
	if categoryName == "" {
		http.Error(w, "Nom de catégorie requis", http.StatusBadRequest)
		return
	}

	// TODO: Ajouter la méthode AddCategory au service
	// err := h.categoryService.AddCategory(categoryName)

	http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
}

func (h *AdminHandler) DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	categoryID := r.FormValue("category_id")
	if categoryID == "" {
		http.Error(w, "ID de catégorie requis", http.StatusBadRequest)
		return
	}

	// TODO: Ajouter la méthode DeleteCategory au service
	// err := h.categoryService.DeleteCategory(categoryID)

	http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
}

func (h *AdminHandler) ManageReportsHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := h.sessionService.GetCurrentUser(r)

	reports, err := h.reportService.GetPendingReports()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des signalements", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"User":      user,
		"Reports":   reports,
		"DarkMode":  getDarkModeFromCookie(r),
		"PageTitle": "Gestion des signalements",
	}

	if err := h.templateService.RenderTemplate(w, "admin/reports.html", data); err != nil {
		http.Error(w, "Erreur de rendu du template: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *AdminHandler) ResolveReportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	reportIDStr := r.FormValue("report_id")
	statusStr := r.FormValue("status")
	resolution := r.FormValue("resolution")

	reportID, err := strconv.Atoi(reportIDStr)
	if err != nil {
		http.Error(w, "ID de signalement invalide", http.StatusBadRequest)
		return
	}

	// Récupérer l'utilisateur actuel
	user, err := h.sessionService.GetCurrentUser(r)
	if err != nil || user == nil {
		http.Error(w, "Utilisateur non identifié", http.StatusUnauthorized)
		return
	}

	// Vérifier que l'état est valide
	var status domain.ReportStatus
	switch statusStr {
	case "resolved":
		status = domain.ReportResolved
	case "dismissed":
		status = domain.ReportDismissed
	default:
		http.Error(w, "État invalide", http.StatusBadRequest)
		return
	}

	adminUser, err := h.userService.GetUserByEmail(user.Email)
	if err != nil {
		log.Printf("Erreur lors de la récupération de l'admin: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	if err := h.reportService.ResolveReport(reportID, adminUser.ID, resolution, status); err != nil {
		http.Error(w, "Erreur lors de la résolution du signalement", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/reports", http.StatusSeeOther)
}
