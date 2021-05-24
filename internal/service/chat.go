package service

import (
	"time"

	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/repository"
)

type ChatsService struct {
	repo repository.Chats
}

func NewChatsService(repo repository.Chats) *ChatsService {
	return &ChatsService{
		repo: repo,
	}
}

func (s *ChatsService) CreateMessage(message *model.Message) (int, error) {
	message.Date = time.Now()
	message.Read = false

	return s.repo.CreateMessage(message)
}
