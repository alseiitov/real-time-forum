package repository

import (
	"database/sql"
)

type ModeatorsRepo struct {
	db *sql.DB
}

func NewModeatorsRepo(db *sql.DB) *ModeatorsRepo {
	return &ModeatorsRepo{db: db}
}
