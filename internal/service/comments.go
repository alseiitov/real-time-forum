package service

import (
	"time"

	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/repository"
	"github.com/alseiitov/real-time-forum/pkg/image"
	uuid "github.com/satori/go.uuid"
)

type CommentsService struct {
	repo repository.Comments
}

func NewCommentsService(repo repository.Comments) *CommentsService {
	return &CommentsService{
		repo: repo,
	}
}

type CreateCommentInput struct {
	UserID int
	PostID int
	Data   string
	Image  string
}

func (s *CommentsService) Create(input CreateCommentInput) (int, error) {
	data, err := image.BytesFromBase64(input.Image)
	if err != nil {
		return 0, err
	}

	err = image.Validate(data)
	if err != nil {
		return 0, err
	}

	newImageName := uuid.NewV4().String() + image.GetExtension(data)
	err = image.Save(data, newImageName)
	if err != nil {
		return 0, err
	}

	comment := model.Comment{
		UserID: input.UserID,
		PostID: input.PostID,
		Data:   input.Data,
		Image:  newImageName,
		Date:   time.Now(),
	}

	id, err := s.repo.Create(comment)
	return id, err
}

func (s *CommentsService) Delete(userID, postID int) error {
	return s.repo.Delete(userID, postID)
}

func (s *CommentsService) GetCommentsByPostID(postID int) ([]model.Comment, error) {
	comments, err := s.repo.GetCommentsByPostID(postID)
	if err != nil {
		return nil, err
	}

	for i := range comments {
		imgBase64, err := image.ReadImage(comments[i].Image)
		if err != nil {
			return nil, err
		}
		comments[i].Image = imgBase64
	}

	return comments, nil
}
