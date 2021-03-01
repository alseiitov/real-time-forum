package service

import (
	"path/filepath"
	"time"

	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/repository"
	"github.com/alseiitov/real-time-forum/pkg/image"
)

const postsForPage = 5

type PostsService struct {
	repo      repository.Posts
	imagesDir string
}

func NewPostsService(repo repository.Posts, imagesDir string) *PostsService {
	return &PostsService{
		repo:      repo,
		imagesDir: imagesDir,
	}
}

type CreatePostInput struct {
	UserID     int
	Title      string
	Data       string
	Image      string
	Categories []model.Categorie
}

func (s *PostsService) Create(input CreatePostInput) (int, error) {
	imageName, err := image.SaveAndGetName(input.Image, s.imagesDir)
	if err != nil {
		return 0, err
	}

	post := model.Post{
		UserID:     input.UserID,
		Title:      input.Title,
		Data:       input.Data,
		Date:       time.Now(),
		Image:      imageName,
		Categories: input.Categories,
	}

	id, err := s.repo.Create(post)
	return id, err
}

func (s *PostsService) GetByID(postID int) (model.Post, error) {
	post, err := s.repo.GetByID(postID)
	if err != nil {
		return post, err
	}

	imgBase64, err := image.ReadImage(filepath.Join(s.imagesDir, post.Image))
	if err != nil {
		return post, err
	}

	post.Image = imgBase64
	return post, nil
}

func (s *PostsService) Delete(userID, postID int) error {
	return s.repo.Delete(userID, postID)
}

func (s *PostsService) GetPostsByCategoryID(categoryID int, page int) ([]model.Post, error) {
	offset := (page - 1) * postsForPage
	return s.repo.GetPostsByCategoryID(categoryID, postsForPage, offset)
}
