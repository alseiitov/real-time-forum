package repository

import (
	"database/sql"

	"github.com/alseiitov/real-time-forum/internal/model"
)

type Users interface {
	Create(user model.User) error
	GetByCredentials(usernameOrEmail, password string) (model.User, error)
	GetByID(userID int) (model.User, error)

	SetSession(session model.Session) error
	DeleteSession(userID int, refreshToken string) error

	CreateModeratorRequest(userID int) error
}

type Moderators interface {
}

type Admins interface {
	GetModeratorRequests() ([]model.ModeratorRequest, error)
	GetModeratorRequestByID(requestID int) (model.ModeratorRequest, error)
	DeleteModeratorRequest(requestID int) error
	UpdateUserRole(userID, role int) error
}

type Categories interface {
	GetAll() ([]model.Category, error)
	GetByID(categoryID int) (model.Category, error)
}

type Posts interface {
	Create(post model.Post) (int, error)
	GetByID(postID int) (model.Post, error)
	Delete(userID int, postID int) error
	GetPostsByCategoryID(categoryID int, limit int, offset int) ([]model.Post, error)
	LikePost(like model.PostLike) error
}

type Comments interface {
	Create(comment model.Comment) (int, error)
	Delete(userID, commentID int) error
	GetCommentsByPostID(postID int, limit int, offset int) ([]model.Comment, error)
	LikeComment(like model.CommentLike) error
}

type Notifications interface {
	Create(notification model.Notification) error
	GetNotifications(userID int) ([]model.Notification, error)
}

type Repositories struct {
	Users         Users
	Moderators    Moderators
	Admins        Admins
	Categories    Categories
	Posts         Posts
	Comments      Comments
	Notifications Notifications
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Users:         NewUserRepo(db),
		Moderators:    NewModeatorsRepo(db),
		Admins:        NewAdminsRepo(db),
		Categories:    NewCategoriesRepo(db),
		Posts:         NewPostsRepo(db),
		Comments:      NewCommentsRepo(db),
		Notifications: NewNotificationsRepo(db),
	}
}
