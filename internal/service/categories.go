package service

import (
	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/repository"
)

type CategoriesService struct {
	repo         repository.Categories
	postsService Posts
}

func NewCategoriesService(repo repository.Categories, postsService Posts) *CategoriesService {
	return &CategoriesService{
		repo:         repo,
		postsService: postsService,
	}
}

func (s *CategoriesService) GetAll() ([]model.Category, error) {
	return s.repo.GetAll()
}

func (s *CategoriesService) GetByID(categoryID int, page int) (model.Category, error) {
	var category model.Category

	category, err := s.repo.GetByID(categoryID)
	if err != nil {
		return category, err
	}

	category.Posts, err = s.postsService.GetPostsByCategoryID(categoryID, page)

	return category, err
}
