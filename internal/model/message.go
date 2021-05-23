package model

import "time"

type Message struct {
	ID          int       `json:"id,omitempty"`
	SenderID    int       `json:"senderID"`
	RecipientID int       `json:"recipientID"`
	Message     string    `json:"message,omitempty"`
	Date        time.Time `json:"date,omitempty"`
	Read        bool      `json:"read"`
}
