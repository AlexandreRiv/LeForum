package middleware

import (
	"LeForum/internal/auth/session"
	"net/http"
)

func AuthMiddleware(sessionService *session.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Vérifier la présence d'une session valide
			currentUser, err := sessionService.GetCurrentUser(r)
			if err != nil || currentUser == nil {
				http.Redirect(w, r, "/auth", http.StatusSeeOther)
				return
			}

			// Utilisateur authentifié, continuer
			next.ServeHTTP(w, r)
		})
	}
}
