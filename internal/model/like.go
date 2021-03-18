package model

type PostLike struct {
	ID       int
	PostID   int
	UserID   int
	LikeType int
}

type CommentLike struct {
	ID        int
	CommentID int
	UserID    int
	LikeType  int
}

var LikeTypes = struct {
	Like    int
	Dislike int
}{
	Like:    1,
	Dislike: 2,
}
