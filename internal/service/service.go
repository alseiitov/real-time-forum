package service

import (
	"github.com/alseiitov/real-time-forum/internal/repository"
	"github.com/alseiitov/real-time-forum/pkg/auth"
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

type UsersSignInInput struct {
	UsernameOrEmail string
	Password        string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Users interface {
	SignUp(input UsersSignUpInput) error
	SignIn(input UsersSignInInput) (Tokens, error)
}

type Services struct {
	Users Users
}

type ServicesDeps struct {
	Repos        *repository.Repositories
	Hasher       hash.PasswordHasher
	TokenManager auth.TokenManager
}

func NewServices(deps ServicesDeps) *Services {
	return &Services{
		Users: NewUsersService(deps.Repos.Users, deps.Hasher, deps.TokenManager),
	}
}
