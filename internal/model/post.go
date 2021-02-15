package model

import "time"

type Post struct {
	ID         int
	UserID     int
	Title      string
	Data       string
	Date       time.Time
	Image      string
	Categories []int
}
