package service

import (
	"time"

	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/repository"
)

type ChatsService struct {
	repo       repository.Chats
	eventsChan chan *model.WSEvent
}

func NewChatsService(repo repository.Chats, eventsChan chan *model.WSEvent) *ChatsService {
	return &ChatsService{
		repo:       repo,
		eventsChan: eventsChan,
	}
}

func (s *ChatsService) CreateMessage(senderID, recipientID int, message string) error {
	var err error

	msg := &model.Message{
		SenderID:    senderID,
		RecipientID: recipientID,
		Message:     message,
		Date:        time.Now(),
		Read:        false,
	}

	msg.ID, err = s.repo.CreateMessage(msg)
	if err != nil {
		return err
	}

	s.eventsChan <- &model.WSEvent{
		Type:        model.WSEventTypes.Message,
		Body:        msg,
		RecipientID: msg.SenderID,
	}

	s.eventsChan <- &model.WSEvent{
		Type:        model.WSEventTypes.Message,
		Body:        msg,
		RecipientID: msg.RecipientID,
	}

	return nil
}

func (s *ChatsService) GetMessages(senderID, recipientID, lastMessageID int) ([]model.Message, error) {
	limit := 10
	return s.repo.GetMessages(senderID, recipientID, lastMessageID, limit)
}

func (s *ChatsService) ReadMessage(recipientID int, messageID int) error {
	message, err := s.repo.ReadMessage(recipientID, messageID)
	if err != nil {
		return err
	}

	s.eventsChan <- &model.WSEvent{
		Type:        model.WSEventTypes.ReadMessageResponse,
		Body:        message.ID,
		RecipientID: message.SenderID,
	}

	return nil
}
