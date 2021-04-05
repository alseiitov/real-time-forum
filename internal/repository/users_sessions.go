package repository

import (
	"github.com/alseiitov/real-time-forum/internal/model"
)

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
		return ErrNoRows
	}
	return err
}
