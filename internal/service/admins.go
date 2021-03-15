package service

import (
	"time"

	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/repository"
)

type AdminsService struct {
	repo                 repository.Admins
	notificationsService Notifications
	usersService         Users
}

func NewAdminsService(repo repository.Admins, notificationsService Notifications, usersService Users) *AdminsService {
	return &AdminsService{
		repo:                 repo,
		notificationsService: notificationsService,
		usersService:         usersService,
	}
}

func (s *AdminsService) GetModeratorRequests() ([]model.ModeratorRequest, error) {
	return s.repo.GetModeratorRequests()
}

func (s *AdminsService) AcceptRequestForModerator(adminID, requestID int) error {
	request, err := s.repo.GetModeratorRequestByID(requestID)
	if err != nil {
		return err
	}

	err = s.UpdateUserRole(request.User.ID, model.Roles.Moderator)
	if err != nil {
		return err
	}

	now := time.Now()

	roleUpdatedNotification := model.Notification{
		RecipientID:  request.User.ID,
		SenderID:     adminID,
		ActivityType: model.NotificationActivities.RoleUpdated,
		Date:         now,
		Status:       model.NotificationStatus.Unread,
	}

	requestAcceptedNotification := model.Notification{
		RecipientID:  request.User.ID,
		SenderID:     adminID,
		ActivityType: model.NotificationActivities.ModeratorRequestAccepted,
		Date:         now,
		Status:       model.NotificationStatus.Unread,
	}

	err = s.notificationsService.Create(roleUpdatedNotification)
	if err != nil {
		return err
	}

	err = s.notificationsService.Create(requestAcceptedNotification)
	if err != nil {
		return err
	}

	return s.DeleteModeratorRequest(requestID)
}

func (s *AdminsService) DeclineRequestForModerator(adminID, requestID int, message string) error {
	request, err := s.repo.GetModeratorRequestByID(requestID)
	if err != nil {
		return err
	}

	requestDeclinedNotification := model.Notification{
		RecipientID:  request.User.ID,
		SenderID:     adminID,
		ActivityType: model.NotificationActivities.ModeratorRequestDeclined,
		Date:         time.Now(),
		Message:      message,
		Status:       model.NotificationStatus.Unread,
	}

	err = s.notificationsService.Create(requestDeclinedNotification)
	if err != nil {
		return err
	}

	return s.DeleteModeratorRequest(requestID)
}

func (s *AdminsService) UpdateUserRole(userID int, role int) error {
	return s.repo.UpdateUserRole(userID, role)
}

func (s *AdminsService) DeleteModeratorRequest(userID int) error {
	return s.repo.DeleteModeratorRequest(userID)
}
