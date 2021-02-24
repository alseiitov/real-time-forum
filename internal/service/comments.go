package service

import (
	"time"

	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/repository"
)

type CommentsService struct {
	repo repository.Comments
}

func NewCommentsService(repo repository.Comments) *CommentsService {
	return &CommentsService{
		repo: repo,
	}
}

func (s *CommentsService) Create(input CreateCommentInput) (int, error) {
	comment := model.Comment{
		UserID: input.UserID,
		PostID: input.PostID,
		Data:   input.Data,
		//TODO: add image upload
		// Image: "",
		Date: time.Now(),
	}

	id, err := s.repo.Create(comment)
	return id, err
}

func (s *CommentsService) Delete(userID, postID int) error {
	return s.repo.Delete(userID, postID)
}

func (s *CommentsService) GetCommentsByPostID(postID int) ([]model.Comment, error) {
	return s.repo.GetCommentsByPostID(postID)
}
