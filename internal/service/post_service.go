package service

import (
	"LeForum/internal/domain"
	"LeForum/internal/storage/repositories"
	"time"
)

type PostService struct {
	repo *repositories.PostRepository
}

func NewPostService(repo *repositories.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) CreatePost(title, content, sessionID, category string) error {
	return s.repo.CreatePost(title, content, sessionID, category, time.Now())
}

func (s *PostService) GetPosts() ([]domain.Post, error) {
	return s.repo.GetPosts()
}

func (s *PostService) LikePost(sessionID string, postID string, likeType int) error {
	return s.repo.LikePost(sessionID, postID, likeType)
}
