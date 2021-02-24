package repository

import (
	"database/sql"

	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/repository/sqlitedb"
)

type Users interface {
	Create(user model.User) error
	GetByCredentials(usernameOrEmail, password string) (model.User, error)
	SetSession(session model.Session) error
	DeleteSession(userID int, refreshToken string) error
	DeleteExpiredSessions() error
}

type Posts interface {
	Create(post model.Post) (int, error)
	GetByID(postID int) (model.Post, error)
	Delete(userID int, postID int) error

	GetCategories() ([]model.Categorie, error)
}

type Comments interface {
	Create(comment model.Comment) (int, error)
	Delete(userID, commentID int) error
	GetCommentsByPostID(postID int) ([]model.Comment, error)
}

type Repositories struct {
	Users    Users
	Posts    Posts
	Comments Comments
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Users:    sqlitedb.NewUserRepo(db),
		Posts:    sqlitedb.NewPostsRepo(db),
		Comments: sqlitedb.NewCommentsRepo(db),
	}
}
