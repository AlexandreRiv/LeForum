package api

import (
	"LeForum/internal/api/handlers"
	"LeForum/internal/api/middleware"
	"LeForum/internal/auth/oauth"
	"LeForum/internal/config"
	"net/http"
)

func SetupRouter(appConfig *config.AppConfig) *http.ServeMux {
	mux := http.NewServeMux()

	// Middlewares
	authMiddleware := middleware.AuthMiddleware(appConfig.SessionService)

	// Fichiers statiques
	fileServer := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Initialisation des handlers OAuth
	githubHandler := oauth.NewGithubHandler(appConfig.UserService, appConfig.SessionService)
	googleHandler := oauth.NewGoogleHandler(appConfig.UserService, appConfig.SessionService)

	// Routes OAuth
	mux.HandleFunc("/auth/github", githubHandler.LoginHandler)
	mux.HandleFunc("/auth/github/callback", githubHandler.CallbackHandler)
	mux.HandleFunc("/auth/google", googleHandler.LoginHandler)
	mux.HandleFunc("/auth/google/callback", googleHandler.CallbackHandler)

	appConfig.AuthHandler.RegisterRoutes(mux)
	appConfig.CategoryHandler.RegisterRoutes(mux)

	// Créer et enregistrer les routes de la page d'accueil
	homeHandler := handlers.NewHomeHandler(
		appConfig.PostService,
		appConfig.CategoryService,
		appConfig.SessionService,
		handlers.NewTemplateService(),
	)
	homeHandler.RegisterRoutes(mux)

	// Création des routes pour les posts
	mux.HandleFunc("/post/create", authMiddleware(http.HandlerFunc(appConfig.PostHandler.CreatePostHandler)).ServeHTTP)
	mux.HandleFunc("/post/like", authMiddleware(http.HandlerFunc(appConfig.PostHandler.LikePostHandler)).ServeHTTP)
	mux.HandleFunc("/post", http.HandlerFunc(appConfig.PostHandler.PostPageHandler))

	// Comm
	mux.HandleFunc("/comment/create", authMiddleware(http.HandlerFunc(appConfig.CommentHandler.CreateCommentHandler)).ServeHTTP)

	mux.HandleFunc("/toggle-theme", middleware.ToggleThemeHandler)

	return mux
}
