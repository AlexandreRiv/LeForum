package api

import (
	"LeForum/internal/api/handlers"
	"LeForum/internal/api/middleware"
	"LeForum/internal/auth/oauth"
	"LeForum/internal/config"
	"LeForum/internal/domain"
	"LeForum/internal/service"
	"LeForum/internal/storage/repositories"
	"net/http"
)

func SetupRouter(config *config.AppConfig, reportService *service.ReportService) http.Handler {
	mux := http.NewServeMux()

	// Middlewares
	authMiddleware := middleware.AuthMiddleware(config.SessionService)
	adminMiddleware := middleware.RoleMiddleware(config.SessionService, config.UserService, domain.RoleAdmin)
	moderatorMiddleware := middleware.RoleMiddleware(config.SessionService, config.UserService, domain.RoleModerator)

	// Fichiers statiques
	fileServer := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// OAuth Handlers
	googleHandler := oauth.NewGoogleHandler(config.UserService, config.SessionService)
	githubHandler := oauth.NewGithubHandler(config.UserService, config.SessionService)

	// Routes publiques
	mux.HandleFunc("/auth/google", googleHandler.LoginHandler)
	mux.HandleFunc("/auth/google/callback", googleHandler.CallbackHandler)
	mux.HandleFunc("/auth/github", githubHandler.LoginHandler)
	mux.HandleFunc("/auth/github/callback", githubHandler.CallbackHandler)
	mux.HandleFunc("/toggle-theme", middleware.ToggleThemeHandler)

	// Enregistrement des routes pour chaque gestionnaire
	config.AuthHandler.RegisterRoutes(mux)
	config.CategoryHandler.RegisterRoutes(mux)

	// Commentaire routes
	mux.HandleFunc("/comment/create", config.CommentHandler.CreateCommentHandler)
	mux.HandleFunc("/comment/like", config.CommentHandler.LikeCommentHandler)

	// Routes protégées par authentification
	mux.Handle("/create-post", authMiddleware(http.HandlerFunc(config.PostHandler.CreatePostHandler)))
	mux.Handle("/edit-post", authMiddleware(http.HandlerFunc(config.PostHandler.EditPostHandler)))
	mux.Handle("/like-post", authMiddleware(http.HandlerFunc(config.PostHandler.LikePostHandler)))

	// Gestionnaire de modération avec routes protégées par rôle de modérateur
	moderationHandler := handlers.NewModerationHandler(
		reportService,
		config.SessionService,
		config.PostService,
		service.NewCommentService(repositories.NewCommentRepository(config.DB.DB)),
		handlers.NewTemplateService(),
	)

	// Routes pour modérateurs
	mux.Handle("/moderation", moderatorMiddleware(http.HandlerFunc(moderationHandler.ModerationDashboard)))
	mux.Handle("/moderation/report", authMiddleware(http.HandlerFunc(moderationHandler.ReportContentHandler))) // Tout utilisateur peut signaler
	mux.Handle("/moderation/delete-post", moderatorMiddleware(http.HandlerFunc(moderationHandler.DeletePostHandler)))
	mux.Handle("/moderation/delete-comment", moderatorMiddleware(http.HandlerFunc(moderationHandler.DeleteCommentHandler)))

	// Routes pour administrateurs
	mux.Handle("/admin", adminMiddleware(http.HandlerFunc(config.AdminHandler.AdminDashboardHandler)))
	mux.Handle("/admin/users", adminMiddleware(http.HandlerFunc(config.AdminHandler.ManageUsersHandler)))
	mux.Handle("/admin/change-role", adminMiddleware(http.HandlerFunc(config.AdminHandler.ChangeUserRoleHandler)))
	mux.Handle("/admin/categories", adminMiddleware(http.HandlerFunc(config.AdminHandler.ManageCategoriesHandler)))
	mux.Handle("/admin/add-category", adminMiddleware(http.HandlerFunc(config.AdminHandler.AddCategoryHandler)))
	mux.Handle("/admin/delete-category", adminMiddleware(http.HandlerFunc(config.AdminHandler.DeleteCategoryHandler)))
	mux.Handle("/admin/reports", adminMiddleware(http.HandlerFunc(config.AdminHandler.ManageReportsHandler)))
	mux.Handle("/admin/resolve-report", adminMiddleware(http.HandlerFunc(config.AdminHandler.ResolveReportHandler)))

	// Appliquer les middlewares globaux
	finalHandler := middleware.ThemeMiddleware(mux)

	return finalHandler
}
