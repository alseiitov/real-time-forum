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
}

type Posts interface {
	Create(post model.Post) (int, error)
	CreateComment(comment model.Comment) (int, error)
}

type Repositories struct {
	Users Users
	Posts Posts
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Users: sqlitedb.NewUserRepo(db),
		Posts: sqlitedb.NewPostsRepo(db),
	}
}
