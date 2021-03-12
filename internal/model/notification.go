package model

import "time"

type Notification struct {
	ID           int
	RecipientID  int
	SenderID     int
	ActivityType int
	ObjectID     int
	Date         time.Time
	Message      string
	Status       int
}

var NotificationActivities = struct {
	PostLiked                int
	PostDisliked             int
	PostCommented            int
	CommentLiked             int
	CommentDisliked          int
	PostModerationApproved   int
	PostModerationDeclined   int
	RoleUpdated              int
	ModeratorRequestAccepted int
	ModeratorRequestDeclined int
}{
	PostLiked:                1,
	PostDisliked:             2,
	PostCommented:            3,
	CommentLiked:             4,
	CommentDisliked:          5,
	PostModerationApproved:   6,
	PostModerationDeclined:   7,
	RoleUpdated:              8,
	ModeratorRequestAccepted: 9,
	ModeratorRequestDeclined: 10,
}

var NotificationStatus = struct {
	Read   int
	Unread int
}{
	Read:   1,
	Unread: 2,
}
