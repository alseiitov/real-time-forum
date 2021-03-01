package service

import (
	"time"

	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/repository"
	"github.com/alseiitov/real-time-forum/pkg/auth"
	"github.com/alseiitov/real-time-forum/pkg/hash"
)

type Users interface {
	SignUp(input UsersSignUpInput) error
	SignIn(input UsersSignInInput) (Tokens, error)
	RefreshTokens(input UsersRefreshTokensInput) (Tokens, error)
	DeleteExpiredSessions()
}

type Posts interface {
	Create(input CreatePostInput) (int, error)
	GetByID(postID int) (model.Post, error)
	Delete(userID int, postID int) error
	GetCategories() ([]model.Categorie, error)
}

type Comments interface {
	Create(input CreateCommentInput) (int, error)
	Delete(userID, postID int) error
	GetCommentsByPostID(postID int) ([]model.Comment, error)
}

type Services struct {
	Users    Users
	Posts    Posts
	Comments Comments
}

type ServicesDeps struct {
	Repos           *repository.Repositories
	Hasher          hash.PasswordHasher
	TokenManager    auth.TokenManager
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	ImagesDir       string
}

func NewServices(deps ServicesDeps) *Services {
	return &Services{
		Users:    NewUsersService(deps.Repos.Users, deps.Hasher, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL),
		Posts:    NewPostsService(deps.Repos.Posts, deps.ImagesDir),
		Comments: NewCommentsService(deps.Repos.Comments, deps.ImagesDir),
	}
}
