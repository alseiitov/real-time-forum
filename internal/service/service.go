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

	CreateModeratorRequest(userID int) error
}

type Moderators interface {
}

type Admins interface {
	GetModeratorRequests() ([]model.ModeratorRequest, error)
	AcceptRequestForModerator(adminID, requestID int) error
	DeclineRequestForModerator(adminID, requestID int, message string) error

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
	LikePost(postID, userID, likeType int) error
}

type Comments interface {
	Create(input CreateCommentInput) (int, error)
	Delete(userID, postID int) error
	GetCommentsByPostID(postID int, page int) ([]model.Comment, error)
	LikeComment(comentID, userID, likeType int) error
}

type Notifications interface {
	GetNotifications(userID int) ([]model.Notification, error)
	Create(notification model.Notification) error
}

type Chats interface {
}

type Websockets interface {
	AddClient(c *client)
	SendToClient(clientID int, event *WSEvent)
}

type Services struct {
	Users         Users
	Moderators    Moderators
	Admins        Admins
	Categories    Categories
	Posts         Posts
	Comments      Comments
	Notifications Notifications
	Chats         Chats
	Websockets    Websockets
}

type ServicesDeps struct {
	Repos                          *repository.Repositories
	Hasher                         hash.PasswordHasher
	TokenManager                   auth.TokenManager
	AccessTokenTTL                 time.Duration
	RefreshTokenTTL                time.Duration
	MaxConnsForUser                int
	MaxMessageSize                 int
	TokenWait                      time.Duration
	WriteWait                      time.Duration
	PongWait                       time.Duration
	PingPeriod                     time.Duration
	ImagesDir                      string
	DefaultMaleAvatar              string
	DefaultFemaleAvatar            string
	PostsForPage                   int
	CommentsForPage                int
	PostsPreModerationIsEnabled    bool
	CommentsPreModerationIsEnabled bool
}

func NewServices(deps ServicesDeps) *Services {
	notificationsService := NewNotificationsService(deps.Repos.Notifications)
	commentsService := NewCommentsService(
		deps.Repos.Comments, notificationsService, deps.CommentsForPage,
		deps.ImagesDir, deps.CommentsPreModerationIsEnabled,
	)
	postsService := NewPostsService(
		deps.Repos.Posts, commentsService, notificationsService, deps.ImagesDir,
		deps.PostsForPage, deps.PostsPreModerationIsEnabled,
	)
	categoriesService := NewCategoriesService(deps.Repos.Categories, postsService)
	usersService := NewUsersService(
		deps.Repos.Users, deps.Hasher, deps.TokenManager, deps.AccessTokenTTL,
		deps.RefreshTokenTTL, deps.ImagesDir, deps.DefaultMaleAvatar, deps.DefaultFemaleAvatar,
	)
	chatsService := NewChatsService(deps.Repos.Chats)
	moderatorsService := NewModeratorsService(deps.Repos.Moderators)
	adminsService := NewAdminsService(deps.Repos.Admins, notificationsService, usersService)
	websocketsService := NewWebsocketsService(
		deps.MaxConnsForUser, deps.MaxMessageSize, deps.TokenWait,
		deps.WriteWait, deps.PongWait, deps.PingPeriod,
	)

	return &Services{
		Users:         usersService,
		Moderators:    moderatorsService,
		Admins:        adminsService,
		Categories:    categoriesService,
		Posts:         postsService,
		Comments:      commentsService,
		Notifications: notificationsService,
		Chats:         chatsService,
		Websockets:    websocketsService,
	}
}
