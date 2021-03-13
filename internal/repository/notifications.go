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

	rows, err := r.db.Query("SELECT * FROM notifications WHERE recipient_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var notification model.Notification
		err = rows.Scan(&notification.ID, &notification.RecipientID, &notification.SenderID, &notification.ActivityType, &notification.ObjectID, &notification.Date, &notification.Message, &notification.Status)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, rows.Err()
}

func (r *NotificationsRepo) Create(n model.Notification) error {
	recipientID, err := r.getRecipientID(n)
	if err != nil {
		return err
	}

	if recipientID == n.SenderID {
		return nil
	}

	stmt, err := r.db.Prepare(`INSERT INTO notifications (recipient_id, sender_id, activity_type, object_id, date, message, status) VALUES ($1, $2, $3, $4, $5, $6, $7)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(&recipientID, &n.SenderID, &n.ActivityType, &n.ObjectID, &n.Date, &n.Message, &n.Status)
	if err != nil {
		return err
	}

	return err
}

func (r *NotificationsRepo) getRecipientID(n model.Notification) (int, error) {
	if n.RecipientID != 0 {
		return n.RecipientID, nil
	}

	var query string
	var id int

	switch n.ActivityType {
	case
		model.NotificationActivities.PostLiked,
		model.NotificationActivities.PostDisliked,
		model.NotificationActivities.PostCommented,
		model.NotificationActivities.PostModerationApproved,
		model.NotificationActivities.PostModerationDeclined:

		query = `SELECT user_id FROM posts WHERE id = $1`
	case
		model.NotificationActivities.CommentLiked,
		model.NotificationActivities.CommentDisliked:

		query = `SELECT user_id FROM comments WHERE id = $1`
	}

	err := r.db.QueryRow(query, n.ObjectID).Scan(&id)
	return id, err
}
