package model

import "time"

type Comment struct {
	ID     int       `json:"id"`
	UserID int       `json:"userID"`
	PostID int       `json:"postID"`
	Data   string    `json:"data"`
	Image  string    `json:"image"`
	Date   time.Time `json:"date"`
}
