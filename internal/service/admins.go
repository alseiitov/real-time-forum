package service

import (
	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/repository"
)

type AdminsService struct {
	repo         repository.Admins
	usersService Users
}

func NewAdminsService(repo repository.Admins, usersService Users) *AdminsService {
	return &AdminsService{
		repo:         repo,
		usersService: usersService,
	}
}

func (s *AdminsService) CreateModeratorRequest(userID int) error {
	return s.repo.CreateModeratorRequest(userID)
}

func (s *AdminsService) DeleteModeratorRequest(userID int) error {
	return s.repo.DeleteModeratorRequest(userID)
}

func (s *AdminsService) GetModeratorRequesters() ([]model.User, error) {
	return s.repo.GetModeratorRequesters()
}

func (s *AdminsService) AcceptRequestForModerator(userID int) error {
	err := s.UpdateUserRole(userID, model.Roles.Moderator)
	if err != nil {
		return err
	}

	return s.DeleteModeratorRequest(userID)

}

func (s *AdminsService) DeclineRequestForModerator(userID int) error {
	return s.DeleteModeratorRequest(userID)
}

func (s *AdminsService) UpdateUserRole(userID int, role int) error {
	return s.repo.UpdateUserRole(userID, role)
}
