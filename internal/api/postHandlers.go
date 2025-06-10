package api

import (
	"LeForum/internal/auth"
	"LeForum/internal/storage"
	"html/template"
	"net/http"
	"strconv"
	"time"
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
		return
	}
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	SQLPostRequest := "INSERT INTO posts (title, content, id_user, created_at) VALUES (?, ?, (SELECT users.id FROM users INNER JOIN sessions ON users.mail = sessions.user_email WHERE sessions.id = ?), ?);"
	SQLAffectRequest := "INSERT INTO affectation VALUES ((SELECT id FROM posts WHERE title = ? AND content = ? AND id_user = (SELECT users.id FROM users INNER JOIN sessions ON users.mail = sessions.user_email WHERE sessions.id = ?)), (SELECT id FROM categories WHERE name = ?));"

	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Cookie 'session_id' non trouvé", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Erreur lors de la lecture du cookie", http.StatusBadRequest)
		return
	}

	createdAt := time.Now().Add(2 * time.Hour)

	_, err = storage.DB.Exec(
		SQLPostRequest,
		r.FormValue("title"),
		r.FormValue("content"),
		cookie.Value,
		createdAt,
	)
	if err != nil {
		return
	}

	_, err = storage.DB.Exec(
		SQLAffectRequest,
		r.FormValue("title"),
		r.FormValue("content"),
		cookie.Value,
		r.FormValue("category"),
	)
	if err != nil {
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	SQLLikeRequest := "INSERT INTO liked_posts VALUES ((SELECT users.id FROM users INNER JOIN sessions ON users.mail = sessions.user_email WHERE sessions.id = ?), ?, ?);"
	SQLDelLikeReq := "DELETE FROM liked_posts WHERE id_user = (SELECT users.id FROM users INNER JOIN sessions ON users.mail = sessions.user_email WHERE sessions.id = ?) AND id_post = ?;"

	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Cookie 'session_id' non trouvé", http.StatusUnauthorized)
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

	_, err = storage.DB.Exec(
		SQLLikeRequest,
		cookie.Value,
		id,
		likeType,
	)
	if err != nil {
		_, err = storage.DB.Exec(
			SQLDelLikeReq,
			cookie.Value,
			id,
		)
		if err != nil {
			http.Error(w, "like error", http.StatusBadRequest)
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}
