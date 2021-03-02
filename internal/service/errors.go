package service

import "errors"

var (
	ErrTooManyCategories = errors.New("too many categories")

	ErrUserNotExist      = errors.New("user doesn't exist")
	ErrUserWrongPassword = errors.New("wrong password or user doesn't exist")
	ErrUserAlreadyExist  = errors.New("user with this email or login already exists")

	ErrSessionNotFound = errors.New("refresh token is expired or already used")

	ErrDeletingPost    = errors.New("post with this id doesn't exist or you have no permissions to delete this post")
	ErrDeletingComment = errors.New("comment with this id doesn't exist or you have no permissions to delete this comment")
)
