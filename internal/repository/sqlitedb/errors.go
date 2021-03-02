package sqlitedb

import (
	"database/sql"
	"errors"
	"strings"
)

var (
	ErrUserWrongPassword = errors.New("wrong password or user doesn't exist")
	ErrUserAlreadyExist  = errors.New("user with this email or login already exists")

	ErrSessionNotFound = errors.New("refresh token is invalid or already used")

	ErrDeletingPost    = errors.New("post with this id doesn't exist or you have no permissions to delete this post")
	ErrDeletingComment = errors.New("comment with this id doesn't exist or you have no permissions to delete this comment")
)

func isNotExistError(err error) bool {
	return err == sql.ErrNoRows
}

func isAlreadyExistError(err error) bool {
	return strings.HasPrefix(err.Error(), "UNIQUE constraint failed:")
}
