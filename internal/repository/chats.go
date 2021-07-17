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
			users.id, 
			users.first_name, 
			users.last_name, 
			users.avatar,
			IFNULL(MAX(m.id), 0) AS last_message_id, 
			IFNULL(m.sender_id, 0) AS last_message_sender_id,
			IFNULL(m.recipient_id, 0) AS last_message_recipent_id,
			IFNULL(m.message, "") AS last_message,
			IFNULL(m.date, 0)  AS last_message_date,
			SUM(CASE WHEN (m.read = 0 AND m.recipient_id = $1) THEN 1 ELSE 0 END) AS unread_messages_count
		FROM
			users
		LEFT JOIN messages m ON 
			m.sender_id = $1 AND m.recipient_id = users.id 
		OR 
			m.sender_id = users.id AND m.recipient_id = $1
		WHERE NOT users.id = $1
		GROUP BY users.id
		ORDER BY 
			m.id DESC, 
			users.first_name ASC
	`,
		userID)

	if err != nil {
		return chats, err
	}
	defer rows.Close()

	for rows.Next() {
		var chat model.Chat

		err := rows.Scan(
			&chat.User.ID,
			&chat.User.FirstName,
			&chat.User.LastName,
			&chat.User.Avatar,
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

	err := r.db.QueryRow(
		`SELECT id, sender_id, recipient_id FROM messages WHERE id = $1 AND recipient_id = $2`, messageID, recipientID,
	).Scan(
		&message.ID, &message.SenderID, &message.RecipientID,
	)

	if err != nil {
		return message, err
	}

	_, err = r.db.Exec(`UPDATE messages SET read = true WHERE id = $1`, messageID)
	return message, err
}
