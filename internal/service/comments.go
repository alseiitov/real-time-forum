package service

import (
	"time"

	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/repository"
	"github.com/alseiitov/real-time-forum/pkg/image"
)

type CommentsService struct {
	repo                           repository.Comments
	commentsForPage                int
	imagesDir                      string
	commentsPreModerationIsEnabled bool
}

func NewCommentsService(repo repository.Comments, commentsForPage int, imagesDir string, commentsPreModerationIsEnabled bool) *CommentsService {
	return &CommentsService{
		repo:                           repo,
		commentsForPage:                commentsForPage,
		imagesDir:                      imagesDir,
		commentsPreModerationIsEnabled: commentsPreModerationIsEnabled,
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

	if s.commentsPreModerationIsEnabled {
		comment.Status = model.CommentStatus.Pending
	} else {
		comment.Status = model.CommentStatus.Approved
	}

	id, err := s.repo.Create(comment)
	return id, err
}

func (s *CommentsService) Delete(userID, postID int) error {
	return s.repo.Delete(userID, postID)
}

func (s *CommentsService) GetCommentsByPostID(postID int, page int) ([]model.Comment, error) {
	offset := (page - 1) * s.commentsForPage

	comments, err := s.repo.GetCommentsByPostID(postID, s.commentsForPage, offset)
	if err != nil {
		return nil, err
	}

	for i := range comments {
		imgBase64, err := image.ReadImage(s.imagesDir, comments[i].Image)
		if err != nil {
			return nil, err
		}
		comments[i].Image = imgBase64
	}

	return comments, nil
}
