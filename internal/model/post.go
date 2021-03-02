package model

import "time"

type Post struct {
	ID         int        `json:"id,omitempty"`
	UserID     int        `json:"userID,omitempty"`
	Title      string     `json:"title,omitempty"`
	Data       string     `json:"data,omitempty"`
	Date       time.Time  `json:"date,omitempty"`
	Image      string     `json:"image,omitempty"`
	Categories []Category `json:"categories,omitempty"`
	Comments   []Comment  `json:"comments,omitempty"`
}
