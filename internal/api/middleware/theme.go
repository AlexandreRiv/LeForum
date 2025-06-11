package middleware

import (
	"context"
	"net/http"
	"time"
)

func getDarkModeFromCookie(r *http.Request) bool {
	cookie, err := r.Cookie("darkMode")
	if err == nil && cookie.Value == "true" {
		return true
	}
	return false
}

func ThemeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		darkMode := getDarkModeFromCookie(r)
		ctx := context.WithValue(r.Context(), "darkMode", darkMode)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ToggleThemeHandler(w http.ResponseWriter, r *http.Request) {
	darkMode := getDarkModeFromCookie(r)
	darkMode = !darkMode

	expiration := time.Now().Add(365 * 24 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:    "darkMode",
		Value:   boolToString(darkMode),
		Expires: expiration,
		Path:    "/",
	})

	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
