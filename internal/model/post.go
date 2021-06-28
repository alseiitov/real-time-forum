package model

import "time"

type Post struct {
	ID         int        `json:"id,omitempty"`
	Status     int        `json:"status,omitempty"`
	Author     User       `json:"author,omitempty"`
	Title      string     `json:"title,omitempty"`
	Data       string     `json:"data,omitempty"`
	Date       time.Time  `json:"date,omitempty"`
	Image      string     `json:"image,omitempty"`
	Categories []Category `json:"categories,omitempty"`
	Comments   []Comment  `json:"comments,omitempty"`
	Rating     int        `json:"rating"`
	UserRate   int        `json:"userRate"`
}

var PostStatus = struct {
	Approved   int
	Pending    int
	Irrelevant int
	Obscene    int
	Illegal    int
	Insulting  int
}{
	Approved:   1,
	Pending:    2,
	Irrelevant: 3,
	Obscene:    4,
	Illegal:    5,
	Insulting:  6,
}
