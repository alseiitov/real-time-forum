package service

import (
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/repository"
	"github.com/alseiitov/real-time-forum/pkg/image"
)

type PostsService struct {
	repo repository.Posts
}

func NewPostsService(repo repository.Posts) *PostsService {
	return &PostsService{
		repo: repo,
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

	post := model.Post{
		UserID:     input.UserID,
		Title:      input.Title,
		Data:       input.Data,
		Date:       time.Now(),
		Image:      newImageName,
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

	imgBase64, err := image.ReadImage(post.Image)
	if err != nil {
		return post, err
	}

	post.Image = imgBase64
	return post, nil
}

func (s *PostsService) Delete(userID, postID int) error {
	return s.repo.Delete(userID, postID)
}

func (s *PostsService) GetCategories() ([]model.Categorie, error) {
	return s.repo.GetCategories()
}

// func (s *PostsService) CreateCategorie() ([]model.Categorie, error) {
// 	return s.repo.GetCategories()
// }
