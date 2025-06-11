package api

import (
	"LeForum/internal/api/handlers"
	"LeForum/internal/api/middleware"
	"LeForum/internal/auth"
	"net/http"
)

func SetupRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Serveur de fichiers statiques
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes principales
	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/categories", handlers.CategoriesHandler)
	mux.HandleFunc("/post", handlers.PostHandler)
	mux.HandleFunc("/toggle-theme", middleware.ToggleThemeHandler)
	mux.HandleFunc("/createPost", handlers.CreatePostHandler)
	mux.HandleFunc("/likePost", handlers.LikePostHandler)

	// Routes d'authentification
	authHandler := auth.NewHandler()
	authHandler.RegisterRoutes(mux)

	return mux
}
