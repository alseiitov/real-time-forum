package model

import "time"

type Comment struct {
	ID     int
	UserID int
	PostID int
	Data   string
	Image  string
	Date   time.Time
}
