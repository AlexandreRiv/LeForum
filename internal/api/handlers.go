package api

import (
	"html/template"
	"net/http"
	"time"
)

type PageData struct {
	DarkMode bool
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/post_page.html"))
	tmpl.Execute(w, nil)
}

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/categories.html"))
	tmpl.Execute(w, nil)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// check if the dark mode cookie exists
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

func ToggleThemeHandler(w http.ResponseWriter, r *http.Request) {
	// read the cookie to check the current theme
	darkMode := false
	cookie, err := r.Cookie("darkMode")
	if err == nil && cookie.Value == "true" {
		darkMode = true
	}

	// reverse the theme
	darkMode = !darkMode

	// set the cookie to 1 year for the expiration
	expiration := time.Now().Add(365 * 24 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:    "darkMode",
		Value:   boolToString(darkMode),
		Expires: expiration,
		Path:    "/",
	})

	// redirect to the referer page
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
