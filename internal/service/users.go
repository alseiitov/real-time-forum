package service

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/alseiitov/real-time-forum/internal/domain"
	"github.com/alseiitov/real-time-forum/internal/repository"
	"github.com/alseiitov/real-time-forum/pkg/hash"
)

type UsersService struct {
	repo   repository.Users
	hasher hash.PasswordHasher
}

func NewUsersService(repo repository.Users, hasher hash.PasswordHasher) *UsersService {
	return &UsersService{
		repo:   repo,
		hasher: hasher,
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

	//TODO: generate and return tokens
	fmt.Println(user)

	return Tokens{}, nil
}
