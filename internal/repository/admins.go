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

func (r *AdminsRepo) DeleteModeratorRequest(requestID int) error {
	_, err := r.db.Exec(`DELETE FROM moderator_requests WHERE (id = $1)`, requestID)
	if err == sql.ErrNoRows {
		return ErrNoRows
	}

	return err
}

func (r *AdminsRepo) GetModeratorRequests() ([]model.ModeratorRequest, error) {
	var requests []model.ModeratorRequest

	rows, err := r.db.Query(`SELECT id FROM moderator_requests`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var requestID int

		err = rows.Scan(&requestID)
		if err != nil {
			return nil, err
		}

		request, err := r.GetModeratorRequestByID(requestID)
		if err != nil {
			return nil, err
		}

		requests = append(requests, request)
	}

	return requests, nil
}

func (r *AdminsRepo) GetModeratorRequestByID(requestID int) (model.ModeratorRequest, error) {
	var request model.ModeratorRequest
	request.ID = requestID

	row := r.db.QueryRow(`
		SELECT
			id,
			username,
			first_name,
			last_name
		FROM
			users
		WHERE
			(
				id = (
					SELECT
						user_id
					FROM
						moderator_requests
					WHERE
						id = $1
				)
			)
		`,
		requestID,
	)

	err := row.Scan(
		&request.User.ID,
		&request.User.Username,
		&request.User.FirstName,
		&request.User.LastName,
	)

	if err == sql.ErrNoRows {
		return request, ErrNoRows
	}

	return request, err
}

func (r *AdminsRepo) UpdateUserRole(userID, role int) error {
	_, err := r.db.Exec("UPDATE users SET role = ? WHERE id = ?", role, userID)

	return err
}
