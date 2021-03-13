package service

import (
	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/repository"
)

type NotificationsService struct {
	repo repository.Notifications
}

func NewNotificationsService(repo repository.Notifications) *NotificationsService {
	return &NotificationsService{
		repo: repo,
	}
}

func (s *NotificationsService) GetNotifications(userID int) ([]model.Notification, error) {
	return s.repo.GetNotifications(userID)
}

func (s *NotificationsService) Create(notification model.Notification) error {
	return s.repo.Create(notification)
}
