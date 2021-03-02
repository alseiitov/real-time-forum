package model

import "time"

type User struct {
	ID         int       `json:"id,omitempty"`
	Username   string    `json:"username,omitempty"`
	FirstName  string    `json:"firstName,omitempty"`
	LastName   string    `json:"lastName,omitempty"`
	Age        int       `json:"age,omitempty"`
	Gender     int       `json:"gender,omitempty"`
	Email      string    `json:"email,omitempty"`
	Password   string    `json:"password,omitempty"`
	Registered time.Time `json:"registered"`
	Role       int       `json:"role,omitempty"`
	Avatar     string    `json:"avatar,omitempty"`
}

type roles struct {
	Guest         int
	User          int
	Moderator     int
	Administrator int
}

var Roles = roles{
	Guest:         1,
	User:          2,
	Moderator:     3,
	Administrator: 4,
}

type gender struct {
	Male   int
	Female int
}

var Gender = gender{
	Male:   1,
	Female: 2,
}
