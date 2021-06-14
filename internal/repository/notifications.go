package repository

import (
	"database/sql"

	"github.com/alseiitov/real-time-forum/internal/model"
)

type NotificationsRepo struct {
	db *sql.DB
}

func NewNotificationsRepo(db *sql.DB) *NotificationsRepo {
	return &NotificationsRepo{db: db}
}

func (r *NotificationsRepo) GetNotifications(userID int) ([]model.Notification, error) {
	var notifications []model.Notification

	rows, err := r.db.Query(`SELECT * FROM notifications WHERE recipient_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var notification model.Notification

		err = rows.Scan(
			&notification.ID,
			&notification.RecipientID,
			&notification.SenderID,
			&notification.ActivityType,
			&notification.ObjectID,
			&notification.Date,
			&notification.Message,
			&notification.Read,
		)

		if err != nil {
			return nil, err
		}

		notifications = append(notifications, notification)
	}

	return notifications, rows.Err()
}

func (r *NotificationsRepo) Create(n model.Notification) error {
	stmt, err := r.db.Prepare(`
		INSERT INTO 
			notifications 
				(recipient_id, sender_id, activity_type, object_id, date, message, read) 
		VALUES 
			($1, $2, $3, $4, $5, $6, $7)`,
	)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		&n.RecipientID,
		&n.SenderID,
		&n.ActivityType,
		&n.ObjectID,
		&n.Date,
		&n.Message,
		&n.Read,
	)

	if err != nil {
		return err
	}

	return err
}
