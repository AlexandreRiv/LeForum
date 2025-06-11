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

	// Enregistrer les routes des handlers
	appConfig.AuthHandler.RegisterRoutes(mux)
	appConfig.PostHandler.RegisterRoutes(mux)
	appConfig.CategoryHandler.RegisterRoutes(mux)

	// Créer et enregistrer les routes de la page d'accueil
	homeHandler := handlers.NewHomeHandler(
		appConfig.PostService,
		appConfig.CategoryService,
		appConfig.SessionService,
		handlers.NewTemplateService(),
	)
	homeHandler.RegisterRoutes(mux)

	// Routes protégées par l'authentification
	postPageMux := http.NewServeMux()
	postPageMux.HandleFunc("/post", appConfig.PostHandler.PostPageHandler)
	mux.Handle("/post", authMiddleware(postPageMux))

	return mux
}
