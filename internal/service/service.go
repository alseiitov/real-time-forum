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
	GetByID(userID int) (model.User, error)

	RefreshTokens(input UsersRefreshTokensInput) (Tokens, error)
}

type Moderators interface {
}

type Admins interface {
	CreateModeratorRequest(userID int) error
	DeleteModeratorRequest(userID int) error
	GetModeratorRequesters() ([]model.User, error)

	AcceptRequestForModerator(userID int) error
	DeclineRequestForModerator(userID int) error

	UpdateUserRole(userID int, role int) error
}

type Categories interface {
	GetAll() ([]model.Category, error)
	GetByID(categoryID int, page int) (model.Category, error)
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
	GetCommentsByPostID(postID int, page int) ([]model.Comment, error)
}

type Services struct {
	Users      Users
	Moderators Moderators
	Admins     Admins
	Categories Categories
	Posts      Posts
	Comments   Comments
}

type ServicesDeps struct {
	Repos                          *repository.Repositories
	Hasher                         hash.PasswordHasher
	TokenManager                   auth.TokenManager
	AccessTokenTTL                 time.Duration
	RefreshTokenTTL                time.Duration
	ImagesDir                      string
	DefaultMaleAvatar              string
	DefaultFemaleAvatar            string
	PostsForPage                   int
	CommentsForPage                int
	PostsPreModerationIsEnabled    bool
	CommentsPreModerationIsEnabled bool
}

func NewServices(deps ServicesDeps) *Services {
	commentsService := NewCommentsService(deps.Repos.Comments, deps.CommentsForPage, deps.ImagesDir, deps.CommentsPreModerationIsEnabled)
	postsService := NewPostsService(deps.Repos.Posts, commentsService, deps.ImagesDir, deps.PostsForPage, deps.PostsPreModerationIsEnabled)
	categoriesService := NewCategoriesService(deps.Repos.Categories, postsService)
	usersService := NewUsersService(deps.Repos.Users, deps.Hasher, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.ImagesDir, deps.DefaultMaleAvatar, deps.DefaultFemaleAvatar)
	moderatorsService := NewModeratorsService(deps.Repos.Moderators)
	adminsService := NewAdminsService(deps.Repos.Admins, usersService)

	return &Services{
		Users:      usersService,
		Moderators: moderatorsService,
		Admins:     adminsService,
		Categories: categoriesService,
		Posts:      postsService,
		Comments:   commentsService,
	}
}
