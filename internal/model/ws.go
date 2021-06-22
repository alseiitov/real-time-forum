package model

type WSEvent struct {
	Type        string      `json:"type"`
	Body        interface{} `json:"body,omitempty"`
	RecipientID int         `json:"recipientID,omitempty"`
}

var WSEventTypes = struct {
	Message             string
	ChatsRequest        string
	ChatsResponse       string
	MessagesRequest     string
	MessagesResponse    string
	ReadMessageRequest  string
	ReadMessageResponse string
	Notification        string
	OnlineUsersRequst   string
	OnlineUsersResponse string
	Error               string
	PingMessage         string
	PongMessage         string
}{
	Message:             "message",
	ChatsRequest:        "chatsRequest",
	ChatsResponse:       "chatsResponse",
	MessagesRequest:     "messagesRequest",
	MessagesResponse:    "messagesResponse",
	ReadMessageRequest:  "readMessageRequest",
	ReadMessageResponse: "readMessageResponse",
	Notification:        "notification",
	OnlineUsersRequst:   "onlineUsersRequest",
	OnlineUsersResponse: "onlineUsersResponse",
	Error:               "error",
	PingMessage:         "pingMessage",
	PongMessage:         "pongMessage",
}
