package domain

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
	Guest:         0,
	User:          1,
	Moderator:     2,
	Administrator: 3,
}

type gender struct {
	Male   int
	Female int
}

var Gender = gender{
	Male:   0,
	Female: 1,
}
