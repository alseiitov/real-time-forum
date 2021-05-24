package handler

import (
	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/validator"
)

type messageInput struct {
	RecipientID int    `json:"recipientID" validator:"required"`
	Message     string `json:"message,omitempty" validator:"required"`
}

func (h *Handler) messageHandler(clientID int, event *WSEvent) error {
	var input messageInput

	err := event.readBodyTo(&input)
	if err != nil {
		return err
	}

	err = validator.Validate(input)
	if err != nil {
		return err
	}

	message := &model.Message{
		SenderID:    clientID,
		RecipientID: input.RecipientID,
		Message:     input.Message,
	}

	message.ID, err = h.chatsService.CreateMessage(message)
	if err != nil {
		return err
	}

	msgEvent := &WSEvent{
		Type: WSEventTypes.Message,
		Body: message,
	}

	sendEventToClient(message.SenderID, msgEvent)    // Send event to message sender
	sendEventToClient(message.RecipientID, msgEvent) // Send event to recipient

	return nil
}
