package service

import (
	"time"

	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/repository"
	"github.com/alseiitov/real-time-forum/pkg/image"
)

type PostsService struct {
	repo                        repository.Posts
	commentsService             Comments
	notificationsService        Notifications
	imagesDir                   string
	postsForPage                int
	postsPreModerationIsEnabled bool
}

func NewPostsService(repo repository.Posts, commentsService Comments, notificationsService Notifications, imagesDir string, postsForPage int, postsPreModerationIsEnabled bool) *PostsService {
	return &PostsService{
		repo:                        repo,
		commentsService:             commentsService,
		notificationsService:        notificationsService,
		imagesDir:                   imagesDir,
		postsForPage:                postsForPage,
		postsPreModerationIsEnabled: postsPreModerationIsEnabled,
	}
}

type CreatePostInput struct {
	UserID     int
	Title      string
	Data       string
	Image      string
	Categories []int
}

func (s *PostsService) Create(input CreatePostInput) (int, error) {
	imageName, err := image.SaveAndGetName(input.Image, s.imagesDir)
	if err != nil {
		return 0, err
	}

	categories := model.CategoriesFromInts(input.Categories)
	if len(categories) > 3 {
		return 0, ErrTooManyCategories
	}

	post := model.Post{
		UserID:     input.UserID,
		Title:      input.Title,
		Data:       input.Data,
		Date:       time.Now(),
		Image:      imageName,
		Categories: categories,
	}

	if s.postsPreModerationIsEnabled {
		post.Status = model.PostStatus.Pending
	} else {
		post.Status = model.PostStatus.Approved
	}

	id, err := s.repo.Create(post)
	if err == repository.ErrForeignKeyConstraint {
		return 0, ErrCategoryDoesntExist
	}

	return id, err
}

func (s *PostsService) GetByID(postID int) (model.Post, error) {
	post, err := s.repo.GetByID(postID)
	if err != nil {
		if err == repository.ErrNoRows {
			return post, ErrPostDoesntExist
		}
		return post, err
	}

	post.Image, err = image.ReadImage(s.imagesDir, post.Image)
	if err != nil {
		return post, err
	}

	return post, err
}

func (s *PostsService) Delete(userID, postID int) error {
	err := s.repo.Delete(userID, postID)
	if err == repository.ErrNoRows {
		return ErrDeletingPost
	}

	return err
}

func (s *PostsService) GetPostsByCategoryID(categoryID int, page int) ([]model.Post, error) {
	offset := (page - 1) * s.postsForPage
	return s.repo.GetPostsByCategoryID(categoryID, s.postsForPage, offset)
}

func (s *PostsService) LikePost(postID, userID, likeType int) error {
	like := model.PostLike{
		PostID:   postID,
		UserID:   userID,
		LikeType: likeType,
	}

	if err := s.repo.LikePost(like); err != nil {
		if err == repository.ErrForeignKeyConstraint {
			return ErrPostDoesntExist
		}
		return err
	}

	// TODO: send notification to post author

	return nil
}
