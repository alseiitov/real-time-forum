package model

type WSEvent struct {
	Type        string      `json:"type"`
	Body        interface{} `json:"body,omitempty"`
	RecipientID int         `json:"recipientID,omitempty"`
}

var WSEventTypes = struct {
	Message          string
	MessagesRequest  string
	MessagesResponse string
	Notification     string
	Error            string
	PingMessage      string
	PongMessage      string
}{
	Message:          "message",
	MessagesRequest:  "messagesRequest",
	MessagesResponse: "messagesResponse",
	Notification:     "notification",
	Error:            "error",
	PingMessage:      "pingMessage",
	PongMessage:      "pongMessage",
}
