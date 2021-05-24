package repository

import (
	"database/sql"

	"github.com/alseiitov/real-time-forum/internal/model"
)

type ChatsRepo struct {
	db *sql.DB
}

func NewChatsRepo(db *sql.DB) *ChatsRepo {
	return &ChatsRepo{db: db}
}

func (r *ChatsRepo) CreateMessage(message *model.Message) (int, error) {
	stmt, err := r.db.Prepare(`INSERT INTO messages (sender_id, recipient_id, message, date, status) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(message.SenderID, message.RecipientID, message.Message, message.Date, message.Read)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	return int(id), err
}
