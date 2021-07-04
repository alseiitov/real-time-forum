package model

type Notification struct {
	ID           int         `json:"id,omitempty"`
	RecipientID  int         `json:"recipientID,omitempty"`
	SenderID     int         `json:"senderID,omitempty"`
	ActivityType int         `json:"activityType,omitempty"`
	ObjectID     int         `json:"objectID,omitempty"`
	Date         interface{} `json:"date,omitempty"`
	Message      string      `json:"message,omitempty"`
	Read         bool        `json:"read,omitempty"`
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
