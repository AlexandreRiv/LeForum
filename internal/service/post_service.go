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
	return s.repo.CreatePost(title, content, sessionID, category, time.Now().Add(2*time.Hour))
}

func (s *PostService) GetPosts(order string) ([]domain.Post, error) {
	return s.repo.GetPosts(order string)
}

func (s *PostService) GetPostByID(id int) (domain.Post, error) {
	return s.repo.GetPostByID(id)
}

func (s *PostService) GetCommentsByPostID(id int) ([]domain.Comment, error) {
	return s.repo.GetCommentsByPostID(id)
}

func (s *PostService) LikePost(sessionID string, postID string, likeType int) error {
	return s.repo.LikePost(sessionID, postID, likeType)
}
