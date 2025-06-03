package main

import (
	"LeForum/internal/api"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", api.Handler)
	http.HandleFunc("/toggle-theme", api.ToggleThemeHandler)
	http.HandleFunc("/post", api.PostHandler)
	http.HandleFunc("/categories", api.CategoriesHandler)

	http.ListenAndServe(":3002", nil)
}
