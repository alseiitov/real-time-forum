package service

import (
	"time"

	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/repository"
	"github.com/alseiitov/real-time-forum/pkg/image"
)

type PostsService struct {
	repo            repository.Posts
	commentsService Comments
	imagesDir       string
	postsForPage    int
}

func NewPostsService(repo repository.Posts, commentsService Comments, imagesDir string, postsForPage int) *PostsService {
	return &PostsService{
		repo:            repo,
		commentsService: commentsService,
		imagesDir:       imagesDir,
		postsForPage:    postsForPage,
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

	id, err := s.repo.Create(post)
	return id, err
}

func (s *PostsService) GetByID(postID int) (model.Post, error) {
	post, err := s.repo.GetByID(postID)
	if err != nil {
		return post, err
	}

	post.Comments, err = s.commentsService.GetCommentsByPostID(postID)
	if err != nil {
		return post, err
	}

	post.Image, err = image.ReadImage(s.imagesDir, post.Image)
	if err != nil {
		return post, err
	}

	return post, err
}

func (s *PostsService) Delete(userID, postID int) error {
	return s.repo.Delete(userID, postID)
}

func (s *PostsService) GetPostsByCategoryID(categoryID int, page int) ([]model.Post, error) {
	offset := (page - 1) * s.postsForPage
	return s.repo.GetPostsByCategoryID(categoryID, s.postsForPage, offset)
}
