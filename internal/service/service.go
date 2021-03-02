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
	GetByID(userID int) (model.User, error)
	DeleteExpiredSessions()
}

type Categories interface {
	GetAll() ([]model.Categorie, error)
	GetByID(categoryID int, page int) (model.Categorie, error)
}

type Posts interface {
	Create(input CreatePostInput) (int, error)
	GetByID(postID int) (model.Post, error)
	Delete(userID int, postID int) error
	GetPostsByCategoryID(categoryID int, page int) ([]model.Post, error)
}

type Comments interface {
	Create(input CreateCommentInput) (int, error)
	Delete(userID, postID int) error
	GetCommentsByPostID(postID int) ([]model.Comment, error)
}

type Services struct {
	Users      Users
	Categories Categories
	Posts      Posts
	Comments   Comments
}

type ServicesDeps struct {
	Repos               *repository.Repositories
	Hasher              hash.PasswordHasher
	TokenManager        auth.TokenManager
	AccessTokenTTL      time.Duration
	RefreshTokenTTL     time.Duration
	ImagesDir           string
	DefaultMaleAvatar   string
	DefaultFemaleAvatar string
	PostsForPage        int
}

func NewServices(deps ServicesDeps) *Services {
	return &Services{
		Users:      NewUsersService(deps.Repos.Users, deps.Hasher, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.ImagesDir, deps.DefaultMaleAvatar, deps.DefaultFemaleAvatar),
		Categories: NewCategoriesService(deps.Repos.Categories),
		Posts:      NewPostsService(deps.Repos.Posts, deps.ImagesDir, deps.PostsForPage),
		Comments:   NewCommentsService(deps.Repos.Comments, deps.ImagesDir),
	}
}
