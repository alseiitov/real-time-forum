package service

import (
	"time"

	"github.com/alseiitov/real-time-forum/internal/domain"
	"github.com/alseiitov/real-time-forum/internal/repository"
)

type PostsService struct {
	repo repository.Posts
}

func NewPostsService(repo repository.Posts) *PostsService {
	return &PostsService{
		repo: repo,
	}
}

func (s *PostsService) Create(input CreatePostInput) error {

	post := domain.Post{
		UserID: input.UserID,
		Title:  input.Title,
		Data:   input.Data,
		Date:   time.Now(),
		//TODO: add image upload
		Image:      "",
		Categories: input.Categories,
	}
	return nil
}
