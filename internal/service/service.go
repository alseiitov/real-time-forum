package service

import "github.com/alseiitov/real-time-forum/internal/repository"

type Users interface {
}

type Services struct {
	Users Users
}

type ServicesDeps struct {
	Repos *repository.Repositories
}

func NewServices(deps ServicesDeps) *Services {
	return &Services{
		Users: NewUsersService(deps.Repos.Users),
	}
}
