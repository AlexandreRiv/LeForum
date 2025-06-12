package middleware

import (
	"LeForum/internal/auth/session"
	"LeForum/internal/domain"
	"LeForum/internal/service"
	"net/http"
)

func RoleMiddleware(sessionService *session.Service, userService *service.UserService, requiredRole domain.RoleType) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Vérifier la session
			sessionUser, err := sessionService.GetCurrentUser(r)
			if err != nil || sessionUser == nil {
				http.Redirect(w, r, "/auth", http.StatusSeeOther)
				return
			}

			// Récupérer l'utilisateur complet avec son rôle
			user, err := userService.GetUserByEmail(sessionUser.Email)
			if err != nil {
				http.Error(w, "Erreur de serveur", http.StatusInternalServerError)
				return
			}

			// Vérifier le rôle
			if !hasRequiredRole(user.Role, requiredRole) {
				http.Error(w, "Accès refusé", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// hasRequiredRole vérifie si le rôle de l'utilisateur a les permissions nécessaires
func hasRequiredRole(userRole domain.RoleType, requiredRole domain.RoleType) bool {
	if userRole == domain.RoleAdmin {
		// L'administrateur a tous les droits
		return true
	}

	if userRole == domain.RoleModerator && requiredRole == domain.RoleUser {
		// Le modérateur a les droits d'utilisateur
		return true
	}

	// Sinon, vérifier l'égalité directe
	return userRole == requiredRole
}
