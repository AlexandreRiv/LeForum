package service

import (
	"LeForum/internal/storage/repositories"
	"time"
)

type CommentService struct {
	repo *repositories.CommentRepository
}

func NewCommentService(repo *repositories.CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) CreateComment(content, sessionID string, postID int) error {
	return s.repo.CreateComment(content, sessionID, postID, time.Now().Add(2*time.Hour))
}