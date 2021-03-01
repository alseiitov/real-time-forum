package service

import (
	"path/filepath"
	"time"

	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/repository"
	"github.com/alseiitov/real-time-forum/pkg/image"
)

type CommentsService struct {
	repo      repository.Comments
	imagesDir string
}

func NewCommentsService(repo repository.Comments, imagesDir string) *CommentsService {
	return &CommentsService{
		repo:      repo,
		imagesDir: imagesDir,
	}
}

type CreateCommentInput struct {
	UserID int
	PostID int
	Data   string
	Image  string
}

func (s *CommentsService) Create(input CreateCommentInput) (int, error) {
	imageName, err := image.SaveAndGetName(input.Image, s.imagesDir)
	if err != nil {
		return 0, err
	}

	comment := model.Comment{
		UserID: input.UserID,
		PostID: input.PostID,
		Data:   input.Data,
		Image:  imageName,
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
		imgBase64, err := image.ReadImage(filepath.Join(s.imagesDir, comments[i].Image))
		if err != nil {
			return nil, err
		}
		comments[i].Image = imgBase64
	}

	return comments, nil
}
