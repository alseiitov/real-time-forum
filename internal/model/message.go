package model

type Message struct {
	ID          int         `json:"id,omitempty"`
	SenderID    int         `json:"senderID"`
	RecipientID int         `json:"recipientID"`
	Message     string      `json:"message,omitempty"`
	Date        interface{} `json:"date,omitempty"`
	Read        bool        `json:"read"`
}
