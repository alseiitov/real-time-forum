package model

import "time"

type Message struct {
	ID      int       `json:"id,omitempty"`
	ChatID  int       `json:"chatID,omitempty"`
	Message string    `json:"message,omitempty"`
	Date    time.Time `json:"date,omitempty"`
	UserID  int       `json:"userID,omitempty"`
	Read    bool      `json:"read"`
}
