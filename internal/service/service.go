package service

import (
	"time"

	"github.com/alseiitov/real-time-forum/internal/model"
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

type CreatePostInput struct {
	UserID     int
	Title      string
	Data       string
	Categories []model.Categorie
}

type CreateCommentInput struct {
	UserID int
	PostID int
	Data   string
}

type Posts interface {
	Create(input CreatePostInput) (int, error)
	GetByID(role int, postID int) (model.Post, error)
	Delete(userID int, role int, postID int) error
	CreateComment(input CreateCommentInput) (int, error)
}

type Services struct {
	Users Users
	Posts Posts
}

type ServicesDeps struct {
	Repos           *repository.Repositories
	Hasher          hash.PasswordHasher
	TokenManager    auth.TokenManager
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewServices(deps ServicesDeps) *Services {
	return &Services{
		Users: NewUsersService(deps.Repos.Users, deps.Hasher, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL),
		Posts: NewPostsService(deps.Repos.Posts),
	}
}
