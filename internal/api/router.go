package api

import (
	"LeForum/internal/api/handlers"
	"LeForum/internal/api/middleware"
	"LeForum/internal/auth/oauth"
	"LeForum/internal/config"
	"LeForum/internal/domain"
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
	mux.HandleFunc("/comment/like", authMiddleware(http.HandlerFunc(appConfig.CommentHandler.LikeCommentHandler)).ServeHTTP)

	mux.HandleFunc("/toggle-theme", middleware.ToggleThemeHandler)

	adminProtectedHandler := middleware.RoleMiddleware(
		appConfig.SessionService,
		appConfig.UserService,
		domain.RoleAdmin)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extraction du chemin restant après /admin ou /admin/
path := r.URL.Path
if len(path) >= 6 && path[:6] == "/admin" {
    path = path[6:]
    if path == "" {
        path = "/"
    }
}

// Créer une nouvelle requête avec le chemin modifié
r2 := new(http.Request)
*r2 = *r
r2.URL.Path = path

// Servir via le AdminHandler
appConfig.AdminHandler.ServeHTTP(w, r2)
	}))

	// Enregistrer le handler protégé pour toutes les routes commençant par /admin/
	mux.Handle("/admin/", adminProtectedHandler)
	mux.Handle("/admin", adminProtectedHandler)

	return mux
}
