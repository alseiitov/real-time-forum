package handler

import (
	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/validator"
)

type messageInput struct {
	RecipientID int    `json:"recipientID" validator:"required"`
	Message     string `json:"message,omitempty" validator:"required"`
}

func (h *Handler) messageHandler(clientID int, event *model.WSEvent) error {
	var input messageInput

	err := UnmarshalEventBody(event, &input)
	if err != nil {
		return err
	}

	err = validator.Validate(input)
	if err != nil {
		return err
	}

	return h.chatsService.CreateMessage(clientID, input.RecipientID, input.Message)
}
