package api

import (
	"LeForum/internal/auth"
	"html/template"
	"log"
	"net/http"
	"time"
)

var templates *template.Template

func init() {
	templates = template.Must(template.ParseGlob("web/templates/*.html"))
	template.Must(templates.ParseGlob("web/templates/components/*.html"))
}

// Handler est la fonction principale qui gère la page d'accueil
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Vérifie si l'utilisateur est connecté
	session, err := auth.GetSession(r)
	isLoggedIn := (err == nil && session != nil)

	// Prépare les données pour le template
	data := struct {
		DarkMode    bool
		IsLoggedIn  bool
		CurrentPage string
	}{
		DarkMode:    getDarkModeFromCookie(r),
		IsLoggedIn:  isLoggedIn,
		CurrentPage: "home",
	}

	// Exécute le template
	err = templates.ExecuteTemplate(w, "home_page.html", data)
	if err != nil {
		log.Printf("Erreur lors de l'exécution du template: %v", err)
		http.Error(w, "Erreur Interne du Serveur", http.StatusInternalServerError)
	}
}

// ToggleThemeHandler bascule entre le mode clair et sombre
func ToggleThemeHandler(w http.ResponseWriter, r *http.Request) {
	// Récupère la valeur actuelle du mode sombre
	darkMode := getDarkModeFromCookie(r)

	// Inverse la valeur
	darkMode = !darkMode

	// Crée un nouveau cookie avec la nouvelle valeur
	cookie := &http.Cookie{
		Name:     "dark_mode",
		Value:    boolToString(darkMode),
		Path:     "/",
		HttpOnly: false,
		Secure:   r.TLS != nil,
		Expires:  time.Now().Add(365 * 24 * time.Hour), // 1 an
	}

	http.SetCookie(w, cookie)

	// Redirige vers la page précédente
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}

	http.Redirect(w, r, referer, http.StatusSeeOther)
}

// Fonction utilitaire pour récupérer le mode sombre depuis le cookie
func getDarkModeFromCookie(r *http.Request) bool {
	cookie, err := r.Cookie("dark_mode")
	if err != nil {
		return false
	}
	return cookie.Value == "true"
}

// Fonction utilitaire pour convertir un booléen en chaîne
func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// Fonction pour gérer la page de catégories
func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		DarkMode    bool
		CurrentPage string
	}{
		DarkMode:    getDarkModeFromCookie(r),
		CurrentPage: "categories",
	}

	err := templates.ExecuteTemplate(w, "categories.html", data)
	if err != nil {
		log.Printf("Erreur lors de l'exécution du template: %v", err)
		http.Error(w, "Erreur Interne du Serveur", http.StatusInternalServerError)
	}
}

// Fonction pour gérer la page d'un post individuel
func PostHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		DarkMode    bool
		CurrentPage string
	}{
		DarkMode:    getDarkModeFromCookie(r),
		CurrentPage: "post",
	}

	err := templates.ExecuteTemplate(w, "post_page.html", data)
	if err != nil {
		log.Printf("Erreur lors de l'exécution du template: %v", err)
		http.Error(w, "Erreur Interne du Serveur", http.StatusInternalServerError)
	}
}
