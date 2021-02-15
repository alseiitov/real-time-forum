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
	if err != nil {
		return err
	}

	return nil
}

func (r *UsersRepo) GetUserByLogin(usernameOrEmail string) (model.User, error) {
	var user model.User

	row := r.db.QueryRow("SELECT id, username, first_name, last_name, age, gender, email, role, avatar FROM users WHERE username = $1 OR email = $1", usernameOrEmail)
	err := row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Age, &user.Gender, &user.Email, &user.Role, &user.Avatar)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *UsersRepo) GetPasswordByLogin(usernameOrEmail string) (string, error) {
	var password string

	row := r.db.QueryRow("SELECT password FROM users WHERE username = $1 OR email = $1", usernameOrEmail)
	err := row.Scan(&password)
	if err != nil {
		return "", err
	}

	return password, nil
}

func (r *UsersRepo) SetSession(session model.Session) error {
	stmt, err := r.db.Prepare("INSERT INTO sessions (user_id, refresh_token, expires_at) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(&session.UserID, &session.RefreshToken, &session.ExpiresAt)
	if err != nil {
		return err
	}

	return nil
}
