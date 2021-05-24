package handler

import (
	"encoding/json"
	"fmt"

	"github.com/alseiitov/real-time-forum/internal/model"
)

func (h *Handler) messageHandler(clientID int, event *WSEvent) error {
	var msg model.Message

	bodyBytes, err := json.Marshal(event.Body.(map[string]interface{}))
	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyBytes, &msg)
	if err != nil {
		return err
	}

	msg.ID, err = h.chatsService.CreateMessage(clientID, &msg)
	if err != nil {
		return err
	}

	msgEvent := &WSEvent{Type: WSEventTypes.Message, Body: msg}
	fmt.Println(msg)

	sendEventToClient(msg.SenderID, msgEvent)    // Send event to message sender
	sendEventToClient(msg.RecipientID, msgEvent) // Send event to recipient

	return nil
}
