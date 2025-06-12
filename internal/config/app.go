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

	UserService     *service.UserService
	PostService     *service.PostService
	CategoryService *service.CategoryService
	SessionService  *session.Service

	AuthHandler     *handlers.AuthHandler
	PostHandler     *handlers.PostHandler
	CategoryHandler *handlers.CategoryHandler
	CommentHandler  *handlers.CommentHandler
	AdminHandler    *handlers.AdminHandler
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

	// Initialiser services
	userService := service.NewUserService(userRepo)
	postService := service.NewPostService(postRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	commentService := service.NewCommentService(commentRepo)
	sessionService := session.NewService(db.DB)

	// Initialiser handlers
	templateService := handlers.NewTemplateService()
	authHandler := handlers.NewAuthHandler(userService, sessionService, templateService)
	postHandler := handlers.NewPostHandler(postService, categoryService, sessionService, templateService)
	categoryHandler := handlers.NewCategoryHandler(categoryService, sessionService, templateService)
	commentHandler := handlers.NewCommentHandler(commentService, sessionService, templateService)
	adminHandler := handlers.NewAdminHandler(userService, sessionService, templateService)

	return &AppConfig{
		DB:                 db,
		UserRepository:     userRepo,
		PostRepository:     postRepo,
		CategoryRepository: categoryRepo,
		UserService:        userService,
		PostService:        postService,
		CategoryService:    categoryService,
		SessionService:     sessionService,
		AuthHandler:        authHandler,
		PostHandler:        postHandler,
		CategoryHandler:    categoryHandler,
		CommentHandler:     commentHandler,
		AdminHandler:       adminHandler,
	}, nil
}
