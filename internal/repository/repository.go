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
}

type Repositories struct {
	Users Users
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Users: sqlitedb.NewUserRepo(db),
	}
}
