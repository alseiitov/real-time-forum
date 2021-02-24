package sqlitedb

import (
	"database/sql"
	"errors"

	"github.com/alseiitov/real-time-forum/internal/model"
)

type PostsRepo struct {
	db *sql.DB
}

func NewPostsRepo(db *sql.DB) *PostsRepo {
	return &PostsRepo{db: db}
}

func (r *PostsRepo) Create(post model.Post) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	stmt, err := tx.Prepare("INSERT INTO posts (user_id, title, data, date, image) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(&post.UserID, &post.Title, &post.Data, &post.Date, &post.Image)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, categorie := range post.Categories {
		stmt, err := tx.Prepare("INSERT INTO posts_categories (post_id, categorie_id) VALUES (?, ?)")
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		_, err = stmt.Exec(&id, &categorie.ID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	return int(id), tx.Commit()
}

func (r *PostsRepo) GetByID(postID int) (model.Post, error) {
	var post model.Post

	row := r.db.QueryRow("SELECT * FROM posts WHERE id = $1", postID)
	err := row.Scan(&post.ID, &post.UserID, &post.Title, &post.Data, &post.Date, &post.Image)
	if err != nil {
		return post, err
	}

	post.Categories, err = r.getPostCategories(postID)
	if err != nil {
		return post, err
	}

	return post, nil
}

func (r *PostsRepo) getPostCategories(postID int) ([]model.Categorie, error) {
	var categories []model.Categorie

	rows, err := r.db.Query("SELECT id, name FROM categories WHERE id IN (SELECT categorie_id FROM posts_categories WHERE post_id = $1)", postID)
	if err != nil {
		return categories, err
	}
	defer rows.Close()

	for rows.Next() {
		var categorie model.Categorie
		err = rows.Scan(&categorie.ID, &categorie.Name)
		if err != nil {
			return categories, err
		}
		categories = append(categories, categorie)
	}

	return categories, rows.Err()
}

func (r *PostsRepo) Delete(userID, postID int) error {
	res, err := r.db.Exec("DELETE FROM posts WHERE (id=$1) and (user_id=$2 OR EXISTS (SELECT * FROM users WHERE id=$2 AND role=$3))", postID, userID, model.Roles.Administrator)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if n == 0 {
		return errors.New("post with this id doesn't exist or you have no permissions to delete this post")
	}

	return err
}

func (r *PostsRepo) GetCategories() ([]model.Categorie, error) {
	var categories []model.Categorie

	rows, err := r.db.Query("SELECT * FROM categories")
	if err != nil {
		return categories, err
	}
	defer rows.Close()

	for rows.Next() {
		var categorie model.Categorie
		err = rows.Scan(&categorie.ID, &categorie.Name)
		if err != nil {
			return categories, err
		}
		categories = append(categories, categorie)
	}

	return categories, rows.Err()
}
