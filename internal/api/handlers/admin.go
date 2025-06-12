package handlers

import (
	"LeForum/internal/api/middleware"
	"LeForum/internal/auth/session"
	"LeForum/internal/service"
	"net/http"
	"strconv"
)

type AdminHandler struct {
	userService     *service.UserService
	sessionService  *session.Service
	templateService *TemplateService
}

func NewAdminHandler(us *service.UserService, ss *session.Service, ts *TemplateService) *AdminHandler {
	return &AdminHandler{
		userService:     us,
		sessionService:  ss,
		templateService: ts,
	}
}

func (h *AdminHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/admin", h.AdminDashboard)
	mux.HandleFunc("/admin/users", h.ManageUsers)
	mux.HandleFunc("/admin/change-role", h.ChangeUserRole)
}

func (h *AdminHandler) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	user, err := h.sessionService.GetCurrentUser(r)
	if err != nil || user == nil {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	data := map[string]interface{}{
		"User":        user,
		"DarkMode":    getDarkModeFromCookie(r),
		"CurrentPage": "admin",
	}

	err = h.templateService.RenderTemplate(w, "admin/dashboard.html", data)
	if err != nil {
		http.Error(w, "Erreur de serveur", http.StatusInternalServerError)
	}
}

func (h *AdminHandler) ManageUsers(w http.ResponseWriter, r *http.Request) {
	// Vérifier que l'utilisateur est connecté et est admin
	currentUser, err := h.sessionService.GetCurrentUser(r)
	if err != nil || currentUser == nil {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	// Vérifier que l'utilisateur a le rôle d'administrateur
	if middleware.RoleType(currentUser.Role) != middleware.RoleAdmin {
		http.Error(w, "Accès refusé", http.StatusForbidden)
		return
	}

	users, err := h.userService.GetAllUsers()
	if err != nil {
		http.Error(w, "Erreur de récupération des utilisateurs", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Users":       users,
		"DarkMode":    getDarkModeFromCookie(r),
		"CurrentPage": "admin",
		"User":        currentUser, // Ajout de l'utilisateur actuel pour l'affichage
	}

	err = h.templateService.RenderTemplate(w, "admin/users.html", data)
	if err != nil {
		http.Error(w, "Erreur de serveur", http.StatusInternalServerError)
	}
}

func (h *AdminHandler) ChangeUserRole(w http.ResponseWriter, r *http.Request) {
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

	if role != string(middleware.RoleUser) &&
		role != string(middleware.RoleModerator) &&
		role != string(middleware.RoleAdmin) {
		http.Error(w, "Rôle invalide", http.StatusBadRequest)
		return
	}

	err = h.userService.UpdateUserRole(userID, middleware.RoleType(role))
	if err != nil {
		http.Error(w, "Erreur lors de la mise à jour du rôle", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

func (h *AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "":
		h.AdminDashboard(w, r)
	case "users":
		h.ManageUsers(w, r)
	case "change-role":
		h.ChangeUserRole(w, r)
	default:
		http.NotFound(w, r)
	}
}
