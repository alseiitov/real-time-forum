package service

import "errors"

var (
	ErrTooManyCategories            = errors.New("too many categories")
	ErrUserDoesNotExist             = errors.New("user doesn't exist")
	ErrPostDoesntExist              = errors.New("post doesn't exist")
	ErrCommentDoesntExist           = errors.New("comment doesn't exist")
	ErrCategoryDoesntExist          = errors.New("category doesn't exist")
	ErrUserWrongPassword            = errors.New("wrong password or username")
	ErrUserAlreadyExist             = errors.New("user with this email or login already exists")
	ErrModeratorRequestAlreadyExist = errors.New("you've already requested moderator role")
	ErrModeratorRequestDoesntExist  = errors.New("moderator request doesn't exist'")
	ErrSessionNotFound              = errors.New("refresh token is expired or already used")
	ErrDeletingPost                 = errors.New("post with this id doesn't exist or you have no permissions to delete this post")
	ErrDeletingComment              = errors.New("comment with this id doesn't exist or you have no permissions to delete this comment")
)
