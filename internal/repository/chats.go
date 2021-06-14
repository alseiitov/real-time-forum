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
	stmt, err := r.db.Prepare(`INSERT INTO messages (sender_id, recipient_id, message, date, read) VALUES (?, ?, ?, ?, ?)`)
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

func (r *ChatsRepo) GetMessages(senderID, recipientID, lastMessageID, limit int) ([]model.Message, error) {
	var messages []model.Message

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(`
		SELECT 
			id, sender_id, recipient_id, message, date, read 
		FROM 
			messages 
		WHERE 
			(
				(sender_id = $1 AND recipient_id = $2) 
				OR 
				(sender_id = $2 AND recipient_id = $1)
			) 
		AND 
			CASE WHEN 
				$3 = 0 
			THEN 
				true 
			ELSE 
				id < $3 
			END
		ORDER BY 
			id 
		DESC LIMIT $4;
	`,
		senderID, recipientID, lastMessageID, limit,
	)

	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var message model.Message

		err := rows.Scan(
			&message.ID,
			&message.SenderID,
			&message.RecipientID,
			&message.Message,
			&message.Date,
			&message.Read,
		)

		if err != nil {
			tx.Rollback()
			return nil, err
		}

		messages = append(messages, message)

		// update message read status to true if user is message recipient
		if message.RecipientID == senderID {
			_, err = tx.Exec(`UPDATE messages SET read = true WHERE id = $1`, message.ID)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	return messages, tx.Commit()
}
