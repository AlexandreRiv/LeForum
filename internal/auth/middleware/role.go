package middleware

import (
	"LeForum/internal/auth/session"
	"net/http"
)

// RoleType définit un type pour les rôles utilisateur
type RoleType string

const (
	RoleUser  RoleType = "user"
	RoleAdmin RoleType = "admin"
)

// RoleMiddleware vérifie si l'utilisateur a le rôle requis
func RoleMiddleware(sessionService *session.Service, requiredRole RoleType) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Vérifier la session
			currentUser, err := sessionService.GetCurrentUser(r)
			if err != nil || currentUser == nil {
				http.Redirect(w, r, "/auth", http.StatusSeeOther)
				return
			}

			// Vous pourriez implémenter ici une vérification de rôle basée sur la structure User
			// Par exemple, en ajoutant un champ Role à la structure User ou en vérifiant en BDD

			// Pour l'instant, tous les utilisateurs authentifiés sont autorisés
			next.ServeHTTP(w, r)
		})
	}
}
