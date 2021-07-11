package model

type Chat struct {
	//  the user with whom chats list requester is communicating
	User                User    `json:"user,omitempty"`
	LastMessage         Message `json:"lastMessage"`
	UnreadMessagesCount int     `json:"unreadMessagesCount"`
}
