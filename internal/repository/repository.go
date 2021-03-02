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
	DeleteExpiredSessions() error
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
}

type Comments interface {
	Create(comment model.Comment) (int, error)
	Delete(userID, commentID int) error
	GetCommentsByPostID(postID int) ([]model.Comment, error)
}

type Repositories struct {
	Users      Users
	Categories Categories
	Posts      Posts
	Comments   Comments
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Users:      NewUserRepo(db),
		Categories: NewCategoriesRepo(db),
		Posts:      NewPostsRepo(db),
		Comments:   NewCommentsRepo(db),
	}
}
