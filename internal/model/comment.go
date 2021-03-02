package model

import "time"

type Comment struct {
	ID     int       `json:"id,omitempty"`
	UserID int       `json:"userID,omitempty"`
	PostID int       `json:"postID,omitempty"`
	Data   string    `json:"data,omitempty"`
	Image  string    `json:"image,omitempty"`
	Date   time.Time `json:"date,omitempty"`
}
