package sqlitedb

import (
	"database/sql"
	"errors"

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

func (r *UsersRepo) SetSession(session model.Session) error {
	stmt, err := r.db.Prepare("INSERT INTO sessions (user_id, refresh_token, expires_at) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(&session.UserID, &session.RefreshToken, &session.ExpiresAt)

	return err
}

func (r *UsersRepo) DeleteSession(userID int, refreshToken string) error {
	res, err := r.db.Exec("DELETE FROM sessions WHERE user_id = $1 AND refresh_token = $2", userID, refreshToken)
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return errors.New("refresh token is invalid or already used")
	}
	return err
}
