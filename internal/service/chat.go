package service

import (
	"time"

	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/repository"
)

type ChatsService struct {
	repo       repository.Chats
	userRepo   repository.Users
	eventsChan chan *model.WSEvent
}

func NewChatsService(repo repository.Chats, usersRepo repository.Users, eventsChan chan *model.WSEvent) *ChatsService {
	return &ChatsService{
		repo:       repo,
		userRepo:   usersRepo,
		eventsChan: eventsChan,
	}
}

func (s *ChatsService) GetChats(userID int) ([]model.Chat, error) {
	return s.repo.GetChats(userID)
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
		Body:        message,
		RecipientID: message.SenderID,
	}

	s.eventsChan <- &model.WSEvent{
		Type:        model.WSEventTypes.ReadMessageResponse,
		Body:        message,
		RecipientID: recipientID,
	}

	return nil
}
