package api

import (
	"LeForum/internal/storage"
	"html/template"
	"net/http"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/post_page.html"))	

	tmpl.Execute(w, nil)
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	SQLRequest := "INSERT INTO posts (title, content, id_user) VALUES (?, ?, (SELECT users.id FROM users INNER JOIN sessions ON users.mail = sessions.user_email WHERE sessions.id = ?));"

	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Cookie 'session_id' non trouvé", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Erreur lors de la lecture du cookie", http.StatusBadRequest)
		return
	}

	_, err = storage.DB.Exec(
		SQLRequest,
		r.FormValue("title"),
		r.FormValue("content"),
		cookie.Value,
	)
	if err != nil {
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}
