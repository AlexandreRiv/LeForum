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

func (s *PostService) CreatePost(title, content, sessionID, category string, image []byte) error {
	return s.repo.CreatePost(title, content, sessionID, category, image, time.Now().Add(2*time.Hour))
}

func (s *PostService) GetPosts(order, search string) ([]domain.Post, error) {
	return s.repo.GetPosts(order, search)
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

func (s *PostService) UpdatePost(postID int, title, content, category string) error {
	return s.repo.UpdatePost(postID, title, content, category)
}

func (s *PostService) DeletePost(postID int) error {
	return s.repo.DeletePost(postID)
}
