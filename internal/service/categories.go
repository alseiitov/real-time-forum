package service

import (
	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/repository"
)

type CategoriesService struct {
	repo repository.Categories
}

func NewCategoriesService(repo repository.Categories) *CategoriesService {
	return &CategoriesService{
		repo: repo,
	}
}

func (s *CategoriesService) GetAll() ([]model.Categorie, error) {
	return s.repo.GetAll()
}

func (s *CategoriesService) GetByID(categoryID int, page int) (model.Categorie, error) {
	return s.repo.GetByID(categoryID)
}
