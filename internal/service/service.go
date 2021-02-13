package service

import (
	"github.com/alseiitov/real-time-forum/internal/repository"
	"github.com/alseiitov/real-time-forum/pkg/hash"
)

type UsersSignUpInput struct {
	Username  string
	FirstName string
	LastName  string
	Age       int
	Gender    int
	Email     string
	Password  string
}

type Users interface {
	SignUp(input UsersSignUpInput) error
}

type Services struct {
	Users Users
}

type ServicesDeps struct {
	Repos  *repository.Repositories
	Hasher hash.PasswordHasher
}

func NewServices(deps ServicesDeps) *Services {
	return &Services{
		Users: NewUsersService(deps.Repos.Users, deps.Hasher),
	}
}
