package service

import (
	"database/sql"
	"errors"
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

func (s *UsersService) SignUp(input UsersSignUpInput) error {
	password, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

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

	return s.createSession(user.ID, user.Role)
}

func (s *UsersService) createSession(userID, role int) (Tokens, error) {
	accessToken, err := s.tokenManager.NewJWT(userID, role)
	refreshToken := s.tokenManager.NewRefreshToken()
	if err != nil {
		return Tokens{}, err
	}

	tokens := Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	session := model.Session{
		UserID:       userID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	return tokens, s.repo.SetSession(session)
}
