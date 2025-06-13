package handlers

import (
	"LeForum/internal/auth/session"
	"LeForum/internal/service"
)

type NotificationHandler struct {
	notificationService  *service.NotificationService
	sessionService  *session.Service
	templateService *TemplateService
}

func NewNotificationHandler(cs *service.NotificationService, ss *session.Service, ts *TemplateService) *NotificationHandler {
	return &NotificationHandler{
		notificationService:  cs,
		sessionService:  ss,
		templateService: ts,
	}
}