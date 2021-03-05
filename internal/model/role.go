package model

type roles struct {
	Guest     int
	User      int
	Moderator int
	Admin     int
}

var Roles = roles{
	Guest:     1,
	User:      2,
	Moderator: 3,
	Admin:     4,
}
