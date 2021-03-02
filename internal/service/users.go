package service

import (
	"time"

	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/repository"
	"github.com/alseiitov/real-time-forum/pkg/auth"
	"github.com/alseiitov/real-time-forum/pkg/hash"
)

type UsersService struct {
	repo            repository.Users
	hasher          hash.PasswordHasher
	tokenManager    auth.TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewUsersService(repo repository.Users, hasher hash.PasswordHasher, tokenManager auth.TokenManager, accessTokenTTL time.Duration, refreshTokenTTL time.Duration) *UsersService {
	return &UsersService{
		repo:            repo,
		hasher:          hasher,
		tokenManager:    tokenManager,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

type UsersSignUpInput struct {
	Username  string
	FirstName string
	LastName  string
	Age       int
	Gender    int
	Email     string
	Password  string
}

func (s *UsersService) SignUp(input UsersSignUpInput) error {
	password := s.hasher.Hash(input.Password)

	user := model.User{
		Username:   input.Username,
		FirstName:  input.FirstName,
		LastName:   input.LastName,
		Age:        input.Age,
		Gender:     input.Gender,
		Email:      input.Email,
		Password:   password,
		Registered: time.Now(),
		Role:       model.Roles.User,
		Avatar:     "",
	}

	err := s.repo.Create(user)
	if err != nil {
		if err == repository.ErrUserAlreadyExist {
			return ErrUserAlreadyExist
		}
		return err
	}

	return nil
}

type UsersSignInInput struct {
	UsernameOrEmail string
	Password        string
}

func (s *UsersService) SignIn(input UsersSignInInput) (Tokens, error) {
	password := s.hasher.Hash(input.Password)

	user, err := s.repo.GetByCredentials(input.UsernameOrEmail, password)
	if err != nil {
		if err == repository.ErrUserWrongPassword {
			return Tokens{}, ErrUserWrongPassword
		}
		return Tokens{}, err
	}

	return s.setSession(user.ID, user.Role)
}

func (s *UsersService) GetByID(userID int) (model.User, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		if err == repository.ErrUserNotExist {
			return user, ErrUserNotExist
		}
		return user, err
	}

	return user, nil
}
