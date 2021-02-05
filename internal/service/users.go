package service

import "github.com/alseiitov/real-time-forum/internal/repository"

type UsersService struct {
	repo repository.Users
}

func NewUsersService(repo repository.Users) *UsersService {
	return &UsersService{
		repo: repo,
	}
}
