package service

import "github.com/alseiitov/real-time-forum/internal/repository"

type ChatsService struct {
	repo repository.Chats
}

func NewChatsService(repo repository.Chats) *ChatsService {
	return &ChatsService{
		repo: repo,
	}
}
