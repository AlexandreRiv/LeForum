package service

import (
	"LeForum/internal/domain"
	"LeForum/internal/storage/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetCategories() ([]string, error) {
	return s.repo.GetCategories()
}

func (s *CategoryService) GetCategoryPosts(categoryName string) ([]domain.Post, error) {
	return s.repo.GetCategoryPosts(categoryName)
}

func (s *CategoryService) GetAllCategoryPosts(categoryName string, limit, offset int) ([]domain.Post, int, error) {
	return s.repo.GetAllCategoryPosts(categoryName, limit, offset)
}
