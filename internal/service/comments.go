package service

import (
	"time"

	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/repository"
	"github.com/alseiitov/real-time-forum/pkg/image"
)

type CommentsService struct {
	repo                           repository.Comments
	notificationsService           Notifications
	commentsForPage                int
	imagesDir                      string
	commentsPreModerationIsEnabled bool
}

func NewCommentsService(repo repository.Comments, notificationsService Notifications, commentsForPage int, imagesDir string, commentsPreModerationIsEnabled bool) *CommentsService {
	return &CommentsService{
		repo:                           repo,
		notificationsService:           notificationsService,
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
	// Create comment
	imageName, err := image.Save(input.Image, s.imagesDir)
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
	if err != nil {
		if err == repository.ErrForeignKeyConstraint {
			return 0, ErrPostDoesntExist
		}
		return 0, err
	}

	// Create notification for post author
	// TODO: fix recipientID
	notification := model.Notification{
		RecipientID:  input.UserID,
		SenderID:     input.UserID,
		ActivityType: model.NotificationActivities.PostCommented,
		ObjectID:     input.PostID,
		Date:         time.Now(),
		Status:       model.NotificationStatus.Unread,
	}

	err = s.notificationsService.Create(notification)
	return id, err
}

func (s *CommentsService) Delete(userID, postID int) error {
	err := s.repo.Delete(userID, postID)
	if err == repository.ErrNoRows {
		return ErrDeletingComment
	}

	return err
}

func (s *CommentsService) GetCommentsByPostID(postID int, page int) ([]model.Comment, error) {
	offset := (page - 1) * s.commentsForPage

	comments, err := s.repo.GetCommentsByPostID(postID, s.commentsForPage, offset)
	if err != nil {
		if err == repository.ErrNoRows {
			return nil, ErrPostDoesntExist
		}
		return nil, err
	}

	return comments, nil
}

func (s *CommentsService) LikeComment(comentID, userID, likeType int) error {

	like := model.CommentLike{
		CommentID: comentID,
		UserID:    userID,
		LikeType:  likeType,
	}

	if err := s.repo.LikeComment(like); err != nil {
		if err == repository.ErrForeignKeyConstraint {
			return ErrCommentDoesntExist
		}
		return err
	}

	// TODO: send notification to comment author

	return nil
}
