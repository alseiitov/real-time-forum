package model

type Post struct {
	ID         int         `json:"id"`
	Status     int         `json:"status,omitempty"`
	Author     User        `json:"author"`
	Title      string      `json:"title"`
	Data       string      `json:"data,omitempty"`
	Date       interface{} `json:"date"`
	Image      string      `json:"image,omitempty"`
	Categories []Category  `json:"categories,omitempty"`
	Comments   []Comment   `json:"comments,omitempty"`
	Rating     int         `json:"rating"`
	UserRate   int         `json:"userRate"`
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
