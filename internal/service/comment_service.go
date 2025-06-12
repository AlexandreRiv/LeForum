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

func (s *CommentService) CreateComment(content, sessionID string, postID int, image []byte) error {
	return s.repo.CreateComment(content, sessionID, postID, image, time.Now().Add(2*time.Hour))
}

func (s *CommentService) LikeComment(sessionID string, commentID string, likeType int) error {
	return s.repo.LikeComment(sessionID, commentID, likeType)
}

func (s *CommentService) DeleteComment(commentID int) error {
	return s.repo.DeleteComment(commentID)
}