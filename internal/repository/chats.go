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

func (r *ChatsRepo) GetChats(userID int) ([]model.Chat, error) {
	var chats []model.Chat

	rows, err := r.db.Query(`
		SELECT 
			MAX(id), sender_id, recipient_id, message, date, SUM(CASE WHEN (read = 0 AND recipient_id = $1) THEN 1 ELSE 0 END) AS unread_count 
		FROM 
			messages 
		WHERE 
			sender_id = $1 
			OR 
			recipient_id = $1 
		GROUP BY 
			MIN(sender_id, recipient_id), MAX(sender_id, recipient_id)
		ORDER BY 
			date 
		DESC
	`,
		userID)

	if err != nil {
		return chats, err
	}
	defer rows.Close()

	for rows.Next() {
		var chat model.Chat

		err := rows.Scan(
			&chat.LastMessage.ID,
			&chat.LastMessage.SenderID,
			&chat.LastMessage.RecipientID,
			&chat.LastMessage.Message,
			&chat.LastMessage.Date,
			&chat.UnreadMessagesCount,
		)

		if err != nil {
			return nil, err
		}

		chats = append(chats, chat)
	}

	return chats, nil
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

	rows, err := r.db.Query(`
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
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (r *ChatsRepo) ReadMessage(recipientID int, messageID int) (model.Message, error) {
	var message model.Message

	tx, err := r.db.Begin()
	if err != nil {
		return message, err
	}

	err = tx.QueryRow(
		`SELECT id, sender_id FROM messages WHERE id = $1 AND recipient_id = $2`, messageID, recipientID,
	).Scan(
		&message.ID, &message.SenderID,
	)

	if err != nil {
		tx.Rollback()
		return message, err
	}

	_, err = tx.Exec(`UPDATE messages SET read = true WHERE id = $1`, messageID)
	if err != nil {
		tx.Rollback()
		return message, err
	}

	return message, tx.Commit()
}
