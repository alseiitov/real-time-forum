package domain

import "time"

type User struct {
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
	User          int
	Moderator     int
	Administrator int
}

var Roles = roles{
	User:          0,
	Moderator:     1,
	Administrator: 2,
}

type gender struct {
	Male   int
	Female int
}

var Gender = gender{
	Male:   0,
	Female: 1,
}
