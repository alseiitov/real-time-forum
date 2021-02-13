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

func (r *UsersRepo) Create(user domain.User) error {
	stmt, err := r.db.Prepare("INSERT INTO users (username, first_name, last_name, age, gender, email, password, role, avatar) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Username, user.FirstName, user.LastName, user.Age, user.Gender, user.Email, user.Password, user.Role, user.Avatar)
	if err != nil {
		return err
	}

	return nil
}

func (r *UsersRepo) GetById(id int) (domain.User, error) {
	var user domain.User

	return user, nil
}
