package sqlitedb

import (
	"database/sql"

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

func (r *PostsRepo) GetByID(postID int, withComments bool) (model.Post, error) {
	var post model.Post
	// Get post
	row := r.db.QueryRow("SELECT * FROM posts WHERE id = $1", postID)
	err := row.Scan(&post.ID, &post.UserID, &post.Title, &post.Data, &post.Date, &post.Image)
	if err != nil {
		return post, err
	}

	// Get post categories
	post.Categories, err = r.getPostCategories(postID)
	if err != nil {
		return post, err
	}

	//Get post comments
	if withComments {
		post.Comments, err = r.getPostComments(postID)
		if err != nil {
			return post, err
		}
	}

	return post, nil
}

func (r *PostsRepo) getPostCategories(postID int) ([]model.Categorie, error) {
	var categories []model.Categorie

	rows, err := r.db.Query("SELECT categorie_id FROM posts_categories WHERE post_id = $1", postID)
	if err != nil {
		return categories, err
	}
	defer rows.Close()

	for rows.Next() {
		var categorie model.Categorie
		err = rows.Scan(&categorie.ID)
		if err != nil {
			return categories, err
		}
		categories = append(categories, categorie)
	}

	return categories, rows.Err()
}

func (r *PostsRepo) getPostComments(postID int) ([]model.Comment, error) {
	var comments []model.Comment

	rows, err := r.db.Query("SELECT * FROM comments WHERE post_id = $1", postID)
	if err != nil {
		return comments, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment model.Comment
		err = rows.Scan(&comment.ID, &comment.UserID, &comment.PostID, &comment.Data, &comment.Image, &comment.Date)
		if err != nil {
			return comments, err
		}
		comments = append(comments, comment)
	}

	return comments, rows.Err()
}

func (r *PostsRepo) Delete(postID int) error {
	_, err := r.db.Exec("DELETE FROM posts WHERE post_id = $1", postID)
	return err
}

func (r *PostsRepo) CreateComment(comment model.Comment) (int, error) {
	stmt, err := r.db.Prepare("INSERT INTO comments (user_id, post_id, data, image, date) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(&comment.UserID, &comment.PostID, &comment.Data, &comment.Image, &comment.Date)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	return int(id), err
}
