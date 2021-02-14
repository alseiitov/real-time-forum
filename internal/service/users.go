package service

import (
	"database/sql"
	"errors"
	"time"

	"github.com/alseiitov/real-time-forum/internal/domain"
	"github.com/alseiitov/real-time-forum/internal/repository"
	"github.com/alseiitov/real-time-forum/pkg/auth"
	"github.com/alseiitov/real-time-forum/pkg/hash"
)

type UsersService struct {
	repo         repository.Users
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager
}

func NewUsersService(repo repository.Users, hasher hash.PasswordHasher, tokenManager auth.TokenManager) *UsersService {
	return &UsersService{
		repo:         repo,
		hasher:       hasher,
		tokenManager: tokenManager,
	}
}

func (s *UsersService) SignUp(input UsersSignUpInput) error {
	password, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	user := domain.User{
		Username:   input.Username,
		FirstName:  input.FirstName,
		LastName:   input.LastName,
		Age:        input.Age,
		Gender:     input.Gender,
		Email:      input.Email,
		Password:   password,
		Registered: time.Now(),
		Role:       domain.Roles.User,
		Avatar:     "",
	}

	err = s.repo.Create(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *UsersService) SignIn(input UsersSignInInput) (Tokens, error) {
	password, err := s.repo.GetPasswordByLogin(input.UsernameOrEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			return Tokens{}, errors.New("Wrong login or password")
		}
		return Tokens{}, err
	}

	if !s.hasher.Compare(input.Password, password) {
		return Tokens{}, errors.New("Wrong login or password")
	}

	user, err := s.repo.GetUserByLogin(input.UsernameOrEmail)
	if err != nil {
		return Tokens{}, err
	}

	accessToken, err := s.tokenManager.NewJWT(user.ID, user.Role)
	if err != nil {
		return Tokens{}, err
	}
	refreshToken := s.tokenManager.NewRefreshToken()

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
