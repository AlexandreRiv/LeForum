package handlers

import (
	"LeForum/internal/auth"
	"LeForum/internal/storage"
	"html/template"
	"net/http"
	"strconv"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetCurrentUser(r)
	darkMode := getDarkModeFromCookie(r)

	data := struct {
		DarkMode    bool
		CurrentPage string
		User        *auth.LoggedUser
	}{
		DarkMode:    darkMode,
		CurrentPage: "post",
		User:        user,
	}

	tmpl := template.Must(template.ParseFiles("web/templates/post_page.html"))
	template.Must(tmpl.ParseGlob("web/templates/components/*.html"))

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Erreur d'affichage du template: "+err.Error(), http.StatusInternalServerError)
	}
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	err = storage.CreatePost(
		r.FormValue("title"),
		r.FormValue("content"),
		cookie.Value,
		r.FormValue("category"),
	)
	if err != nil {
		http.Error(w, "Erreur lors de la création du post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/auth", http.StatusSeeOther)
			return
		}
		http.Error(w, "Erreur lors de la lecture du cookie", http.StatusBadRequest)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Id parameter is missing", http.StatusBadRequest)
		return
	}

	likeType, err := strconv.Atoi(r.URL.Query().Get("like"))
	if err != nil {
		http.Error(w, "Like parameter is missing", http.StatusBadRequest)
		return
	}

	err = storage.LikePost(cookie.Value, id, likeType)
	if err != nil {
		http.Error(w, "like error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
