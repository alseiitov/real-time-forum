package sqlitedb

import (
	"database/sql"

	"github.com/alseiitov/real-time-forum/internal/model"
)

type UsersRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UsersRepo {
	return &UsersRepo{db: db}
}

func (r *UsersRepo) Create(user model.User) error {
	stmt, err := r.db.Prepare("INSERT INTO users (username, first_name, last_name, age, gender, email, password, role, avatar) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(&user.Username, user.FirstName, user.LastName, user.Age, user.Gender, user.Email, user.Password, user.Role, user.Avatar)

	return err
}

func (r *UsersRepo) GetByCredentials(usernameOrEmail, password string) (model.User, error) {
	var user model.User

	row := r.db.QueryRow("SELECT id, role FROM users WHERE (username = $1 OR email = $1) AND (password = $2)", usernameOrEmail, password)
	err := row.Scan(&user.ID, &user.Role)

	return user, err
}
