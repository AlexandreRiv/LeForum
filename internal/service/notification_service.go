package service

import (
	"LeForum/internal/domain"
	"LeForum/internal/storage/repositories"
	"time"
)

type NotificationService struct {
	repo *repositories.NotificationRepository
}

func NewNotificationService(repo *repositories.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) CreateNotification(postID int, content string) error {
	return s.repo.CreateNotification(postID, content, time.Now().Add(2*time.Hour))
}

func (s *NotificationService) GetNotifications(sessionID string) ([]domain.Notification, error) {
	return s.repo.GetNotifications(sessionID)
}