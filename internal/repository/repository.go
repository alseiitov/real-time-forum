package repository

import (
	"database/sql"

	"github.com/alseiitov/real-time-forum/internal/domain"
	"github.com/alseiitov/real-time-forum/internal/repository/sqlitedb"
)

type Users interface {
	Create(user domain.User) error
	GetUserByLogin(usernameOrEmail string) (domain.User, error)
	GetPasswordByLogin(usernameOrEmail string) (string, error)
	SetSession(session domain.Session) error
	// CreateRefreshToken(userID int, refreshToken string, exp )
	// UpdateRefreshToken(userID int, refreshToken string) error
	// DeleteRefreshToken(userID int, refreshToken string) error
}

type Posts interface {
	Create(post domain.Post) (int, error)
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
