package model

import "time"

type Post struct {
	ID         int         `json:"id"`
	UserID     int         `json:"userID"`
	Title      string      `json:"title"`
	Data       string      `json:"data"`
	Date       time.Time   `json:"date"`
	Image      string      `json:"image"`
	Categories []Categorie `json:"categories"`
	Comments   []Comment   `json:"comments"`
}
