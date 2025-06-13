package config

import (
	"LeForum/internal/api/handlers"
	"LeForum/internal/auth/session"
	"LeForum/internal/service"
	"LeForum/internal/storage"
	"LeForum/internal/storage/repositories"
)

type AppConfig struct {
	DB                 *storage.Database
	UserRepository     *repositories.UserRepository
	PostRepository     *repositories.PostRepository
	CategoryRepository *repositories.CategoryRepository
	NotificationRepository *repositories.NotificationRepository

	UserService     *service.UserService
	PostService     *service.PostService
	CategoryService *service.CategoryService
	NotificationService *service.NotificationService
	SessionService  *session.Service

	AuthHandler     *handlers.AuthHandler
	PostHandler     *handlers.PostHandler
	CategoryHandler *handlers.CategoryHandler
	CommentHandler  *handlers.CommentHandler
	NotificationHandler *handlers.NotificationHandler
	AdminHandler    *handlers.AdminHandler

	ReportRepository  *repositories.ReportRepository
	ReportService     *service.ReportService
	ModerationHandler *handlers.ModerationHandler
}

func NewAppConfig() (*AppConfig, error) {
	db, err := storage.InitDB()
	if err != nil {
		return nil, err
	}

	// Initialiser repositories
	userRepo := repositories.NewUserRepository(db.DB)
	postRepo := repositories.NewPostRepository(db.DB)
	categoryRepo := repositories.NewCategoryRepository(db.DB)
	commentRepo := repositories.NewCommentRepository(db.DB)
	notificationRepo := repositories.NewNotificationRepository(db.DB)

	// Initialiser le repository et service pour les rapports
	reportRepo := repositories.NewReportRepository(db.DB)
	reportService := service.NewReportService(reportRepo)

	// Initialiser services
	userService := service.NewUserService(userRepo)
	postService := service.NewPostService(postRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	commentService := service.NewCommentService(commentRepo)
	notificationService := service.NewNotificationService(notificationRepo)
	sessionService := session.NewService(db.DB)

	// Initialiser handlers
	templateService := handlers.NewTemplateService()
	authHandler := handlers.NewAuthHandler(userService, notificationService, sessionService, templateService)
	postHandler := handlers.NewPostHandler(postService, notificationService, categoryService, sessionService, templateService)
	categoryHandler := handlers.NewCategoryHandler(categoryService, notificationService, sessionService, templateService)
	commentHandler := handlers.NewCommentHandler(commentService, notificationService, sessionService, templateService)
	adminHandler := handlers.NewAdminHandler(userService, categoryService, reportService, sessionService, templateService)
	moderationHandler := handlers.NewModerationHandler(reportService, sessionService, postService, commentService, templateService)
	notificationHandler := handlers.NewNotificationHandler(notificationService, sessionService, templateService)

	return &AppConfig{
		DB:                 db,
		UserRepository:     userRepo,
		PostRepository:     postRepo,
		CategoryRepository: categoryRepo,
		NotificationRepository: notificationRepo,
		UserService:        userService,
		PostService:        postService,
		CategoryService:    categoryService,
		NotificationService: notificationService,
		SessionService:     sessionService,
		AuthHandler:        authHandler,
		PostHandler:        postHandler,
		CategoryHandler:    categoryHandler,
		CommentHandler:     commentHandler,
		AdminHandler:       adminHandler,
		ReportRepository:   reportRepo,
		ReportService:      reportService,
		ModerationHandler:  moderationHandler,
		NotificationHandler: notificationHandler,
	}, nil
}
