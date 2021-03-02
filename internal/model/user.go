package model

import "time"

type User struct {
	ID         int
	Username   string
	FirstName  string
	LastName   string
	Age        int
	Gender     int
	Email      string
	Password   string
	Registered time.Time
	Role       int
	Avatar     string
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
