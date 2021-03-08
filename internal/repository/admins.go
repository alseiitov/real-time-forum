package repository

import (
	"database/sql"

	"github.com/alseiitov/real-time-forum/internal/model"
)

type AdminsRepo struct {
	db *sql.DB
}

func NewAdminsRepo(db *sql.DB) *AdminsRepo {
	return &AdminsRepo{db: db}
}

func (r *AdminsRepo) CreateModeratorRequest(userID int) error {
	_, err := r.db.Exec("INSERT INTO moderator_requests (user_id) VALUES ($1)", userID)

	return err
}

func (r *AdminsRepo) DeleteModeratorRequest(userID int) error {
	_, err := r.db.Exec("DELETE FROM moderator_requests WHERE (user_id = $1)", userID)

	return err
}

func (r *AdminsRepo) GetModeratorRequesters() ([]model.User, error) {
	var users []model.User

	rows, err := r.db.Query("SELECT id, username, first_name, last_name FROM users WHERE (id IN (SELECT user_id FROM moderator_requests))")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		err = rows.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *AdminsRepo) UpdateUserRole(userID, role int) error {
	stmt, err := r.db.Prepare("UPDATE users SET role = ? WHERE id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(role, userID)

	return err
}
