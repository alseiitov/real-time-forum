package service

import (
	"time"

	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/repository"
	"github.com/alseiitov/real-time-forum/pkg/image"
)

type CommentsService struct {
	repo                           repository.Comments
	postsRepo                      repository.Posts
	notificationsService           Notifications
	commentsForPage                int
	imagesDir                      string
	commentsPreModerationIsEnabled bool
}

func NewCommentsService(repo repository.Comments, postsRepo repository.Posts, notificationsService Notifications, commentsForPage int,
	imagesDir string, commentsPreModerationIsEnabled bool) *CommentsService {
	return &CommentsService{
		repo:                           repo,
		postsRepo:                      postsRepo,
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
	post, err := s.postsRepo.GetByID(input.PostID, 0)
	if err != nil {
		if err == repository.ErrForeignKeyConstraint {
			return 0, ErrPostDoesntExist
		}
		return 0, err
	}

	notification := model.Notification{
		RecipientID:  post.UserID,
		SenderID:     input.UserID,
		ActivityType: model.NotificationActivities.PostCommented,
		ObjectID:     input.PostID,
		Date:         time.Now(),
		Read:         false,
	}

	return id, s.notificationsService.Create(notification)
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

	likeCreated, err := s.repo.LikeComment(like)
	if err != nil {
		if err == repository.ErrForeignKeyConstraint {
			return ErrCommentDoesntExist
		}
		return err
	}

	// send notification to comment author
	if likeCreated {
		comment, err := s.repo.GetByID(comentID)
		if err != nil {
			return err
		}

		var activityType int

		if likeType == model.LikeTypes.Like {
			activityType = model.NotificationActivities.CommentLiked
		} else {
			activityType = model.NotificationActivities.CommentDisliked
		}

		notification := model.Notification{
			RecipientID:  comment.UserID,
			SenderID:     userID,
			ActivityType: activityType,
			Date:         time.Now(),
			Read:         false,
		}

		return s.notificationsService.Create(notification)
	}

	return nil
}
