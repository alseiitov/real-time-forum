package sqlitedb

import (
	"database/sql"

	"github.com/alseiitov/real-time-forum/internal/domain"
)

type PostsRepo struct {
	db *sql.DB
}

func NewPostsRepo(db *sql.DB) *PostsRepo {
	return &PostsRepo{db: db}
}

func (r *PostsRepo) Create(post domain.Post) error {
	return nil
}
