package ws

import (
	"fmt"

	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/validator"
)

type messageInput struct {
	RecipientID int    `json:"recipientID" validator:"required"`
	Message     string `json:"message,omitempty" validator:"required"`
}

func (h *Handler) newMessage(clientID int, event *model.WSEvent) error {
	var input messageInput

	err := unmarshalEventBody(event, &input)
	if err != nil {
		return err
	}

	err = validator.Validate(input)
	if err != nil {
		return err
	}

	return h.chatsService.CreateMessage(clientID, input.RecipientID, input.Message)
}

type getMessagesInput struct {
	UserID        int `json:"userID" validator:"required"` // id of the user client want to receive messages with
	LastMessageID int `json:"lastMessageID"`
}

func (h *Handler) getMessages(clientID int, event *model.WSEvent) error {
	var input getMessagesInput

	err := unmarshalEventBody(event, &input)
	if err != nil {
		return err
	}

	err = validator.Validate(input)
	if err != nil {
		return err
	}

	messages, err := h.chatsService.GetMessages(clientID, input.UserID, input.LastMessageID)
	if err != nil {
		return err
	}

	h.sendEventToClient(
		&model.WSEvent{
			Type:        model.WSEventTypes.MessagesResponse,
			Body:        messages,
			RecipientID: clientID,
		},
	)

	return nil
}

type readMessageInput struct {
	MessageID int `json:"messageID" validator:"required"`
}

func (h *Handler) readMessage(clientID int, event *model.WSEvent) error {
	var input readMessageInput

	err := unmarshalEventBody(event, &input)
	if err != nil {
		return err
	}

	err = validator.Validate(input)
	if err != nil {
		return err
	}

	err = h.chatsService.ReadMessage(clientID, input.MessageID)
	if err != nil {
		fmt.Println(clientID)
		fmt.Println(input)
		return err
	}

	return nil
}
