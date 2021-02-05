package sqlitedb

import (
	"database/sql"

	"github.com/alseiitov/real-time-forum/internal/domain"
)

type UsersRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UsersRepo {
	return &UsersRepo{db: db}
}

func (r *UsersRepo) GetById(id int) (domain.User, error) {
	var user domain.User
	user.ID = id
	return user, nil
}
