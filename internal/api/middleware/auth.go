package middleware

import (
	"LeForum/internal/auth"
	"context"
	"log"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.GetCurrentUser(r)
		if err != nil {
			log.Printf("Erreur middleware auth: %v", err)
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
