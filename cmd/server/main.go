/*
package main

import (

	"html/template"
	"net/http"

)

	func handler(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("web/templates/home_page.html"))
		tmpl.Execute(w, nil)
	}

	func main() {
		fs := http.FileServer(http.Dir("web/static"))
		http.Handle("/static/", http.StripPrefix("/static/", fs))
		http.HandleFunc("/", handler)

		http.ListenAndServe(":3002", nil)
	}
*/
package main

import (
	"html/template"
	"net/http"
	"time"
)

type PageData struct {
	DarkMode bool
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Vérifier le cookie pour le mode sombre
	darkMode := false
	cookie, err := r.Cookie("darkMode")
	if err == nil && cookie.Value == "true" {
		darkMode = true
	}

	data := PageData{
		DarkMode: darkMode,
	}

	tmpl := template.Must(template.ParseFiles("web/templates/home_page.html"))
	tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Erreur d'affichage du template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func toggleThemeHandler(w http.ResponseWriter, r *http.Request) {
	// Lire l'état actuel du thème
	darkMode := false
	cookie, err := r.Cookie("darkMode")
	if err == nil && cookie.Value == "true" {
		darkMode = true
	}

	// Inverser l'état du thème
	darkMode = !darkMode

	// Définir le cookie avec une expiration d'un an
	expiration := time.Now().Add(365 * 24 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:    "darkMode",
		Value:   boolToString(darkMode),
		Expires: expiration,
		Path:    "/",
	})

	// Rediriger vers la page précédente
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

func main() {
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", handler)
	http.HandleFunc("/toggle-theme", toggleThemeHandler)

	http.ListenAndServe(":3002", nil)
}
