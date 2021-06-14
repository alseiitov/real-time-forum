package service

import (
	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/repository"
)

type NotificationsService struct {
	repo       repository.Notifications
	eventsChan chan *model.WSEvent
}

func NewNotificationsService(repo repository.Notifications, eventsChan chan *model.WSEvent) *NotificationsService {
	return &NotificationsService{
		repo:       repo,
		eventsChan: eventsChan,
	}
}

func (s *NotificationsService) GetNotifications(userID int) ([]model.Notification, error) {
	return s.repo.GetNotifications(userID)
}

func (s *NotificationsService) Create(notification model.Notification) error {
	if notification.RecipientID == notification.SenderID {
		return nil
	}

	err := s.repo.Create(notification)
	if err != nil {
		return err
	}

	s.eventsChan <- &model.WSEvent{
		Type:        model.WSEventTypes.Notification,
		Body:        notification,
		RecipientID: notification.RecipientID,
	}

	return nil
}
