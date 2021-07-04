package model

type Comment struct {
	ID       int         `json:"id,omitempty"`
	Status   int         `json:"status,omitempty"`
	Author   User        `json:"author,omitempty"`
	PostID   int         `json:"postID,omitempty"`
	Data     string      `json:"data,omitempty"`
	Image    string      `json:"image,omitempty"`
	Date     interface{} `json:"date,omitempty"`
	UserRate int         `json:"userRate"`
	Rating   int         `json:"rating"`
}

var CommentStatus = struct {
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
