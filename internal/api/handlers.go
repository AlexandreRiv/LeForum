package api

import (
	"html/template"
	"log"
	"net/http"
)

// Récupère les templates pour le rendu des pages
var templates = template.Must(template.ParseGlob("web/templates/*.html"))

// Handler principal pour la page d'accueil
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	darkMode := getDarkModeFromCookie(r)

	data := struct {
		DarkMode    bool
		CurrentPage string
	}{
		DarkMode:    darkMode,
		CurrentPage: "home",
	}

	err := templates.ExecuteTemplate(w, "home_page.html", data)
	if err != nil {
		log.Printf("Erreur d'affichage du template: %v", err)
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
	}
}

// Récupère la préférence de mode sombre depuis le cookie
func getDarkModeFromCookie(r *http.Request) bool {
	cookie, err := r.Cookie("dark_mode")
	if err != nil {
		// Le cookie n'existe pas, retourner la valeur par défaut (false)
		return false
	}
	return cookie.Value == "true"
}

// Bascule entre le mode clair et le mode sombre
func ToggleThemeHandler(w http.ResponseWriter, r *http.Request) {
	darkMode := !getDarkModeFromCookie(r)

	value := "false"
	if darkMode {
		value = "true"
	}

	cookie := &http.Cookie{
		Name:     "dark_mode",
		Value:    value,
		Path:     "/",
		MaxAge:   365 * 24 * 60 * 60, // 1 an
		HttpOnly: false,
	}

	http.SetCookie(w, cookie)

	// Rediriger vers la page précédente
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}

	http.Redirect(w, r, referer, http.StatusSeeOther)
}
