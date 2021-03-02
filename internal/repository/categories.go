package repository

import (
	"database/sql"

	"github.com/alseiitov/real-time-forum/internal/model"
)

type CategoriesRepo struct {
	db *sql.DB
}

func NewCategoriesRepo(db *sql.DB) *CategoriesRepo {
	return &CategoriesRepo{db: db}
}

func (r *CategoriesRepo) GetAll() ([]model.Category, error) {
	var categories []model.Category

	rows, err := r.db.Query("SELECT * FROM categories")
	if err != nil {
		return categories, err
	}
	defer rows.Close()

	for rows.Next() {
		var category model.Category
		err = rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return categories, err
		}
		categories = append(categories, category)
	}

	return categories, rows.Err()
}

func (r *CategoriesRepo) GetByID(categoryID int) (model.Category, error) {
	var category model.Category

	row := r.db.QueryRow("SELECT id, name FROM categories WHERE id = $1", categoryID)
	err := row.Scan(&category.ID, &category.Name)
	if err != nil {
		return category, err
	}

	return category, nil
}
