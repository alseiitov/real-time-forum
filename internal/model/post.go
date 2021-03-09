package model

import "time"

type Post struct {
	ID         int        `json:"id,omitempty"`
	Status     int        `json:"status,omitempty"`
	UserID     int        `json:"userID,omitempty"`
	Title      string     `json:"title,omitempty"`
	Data       string     `json:"data,omitempty"`
	Date       time.Time  `json:"date,omitempty"`
	Image      string     `json:"image,omitempty"`
	Categories []Category `json:"categories,omitempty"`
	Comments   []Comment  `json:"comments,omitempty"`
}

type postStatus struct {
	Approved   int
	Pending    int
	Irrelevant int
	Obscene    int
	Illegal    int
	Insulting  int
}

var PostStatus = &postStatus{
	Approved:   1,
	Pending:    2,
	Irrelevant: 3,
	Obscene:    4,
	Illegal:    5,
	Insulting:  6,
}
