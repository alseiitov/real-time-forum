package service

import (
	"time"

	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/repository"
	"github.com/alseiitov/real-time-forum/pkg/auth"
	"github.com/alseiitov/real-time-forum/pkg/hash"
	"github.com/alseiitov/real-time-forum/pkg/image"
)

type UsersService struct {
	repo                 repository.Users
	notificationsService Notifications
	hasher               hash.PasswordHasher
	tokenManager         auth.TokenManager
	accessTokenTTL       time.Duration
	refreshTokenTTL      time.Duration
	imagesDir            string
	defaultMaleAvatar    string
	defaultFemaleAvatar  string
}

func NewUsersService(repo repository.Users, hasher hash.PasswordHasher, tokenManager auth.TokenManager, accessTokenTTL time.Duration, refreshTokenTTL time.Duration, imagesDir, defaultMaleAvatar, defaultFemaleAvatar string) *UsersService {
	return &UsersService{
		repo:                repo,
		hasher:              hasher,
		tokenManager:        tokenManager,
		accessTokenTTL:      accessTokenTTL,
		refreshTokenTTL:     refreshTokenTTL,
		imagesDir:           imagesDir,
		defaultMaleAvatar:   defaultMaleAvatar,
		defaultFemaleAvatar: defaultFemaleAvatar,
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
	var avatar string

	if input.Gender == model.Gender.Male {
		avatar = s.defaultMaleAvatar
	}
	if input.Gender == model.Gender.Female {
		avatar = s.defaultFemaleAvatar
	}

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
		Avatar:     avatar,
	}

	err := s.repo.Create(user)

	if err == repository.ErrUserAlreadyExist {
		return ErrUserAlreadyExist
	}

	return err
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

	avatarBase64, err := image.ReadImage(s.imagesDir, user.Avatar)
	if err != nil {
		return user, err
	}

	user.Avatar = avatarBase64

	return user, nil
}

func (s *UsersService) CreateModeratorRequest(userID int) error {
	return s.repo.CreateModeratorRequest(userID)
}
