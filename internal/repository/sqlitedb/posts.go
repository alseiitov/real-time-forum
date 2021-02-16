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
	stmt, err := r.db.Prepare("INSERT INTO posts (user_id, title, data, date, image) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(&post.UserID, &post.Title, &post.Data, &post.Date, &post.Image)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	err = r.addCategoriesToPost(int(id), post.Categories)
	return int(id), err
}

func (r *PostsRepo) addCategoriesToPost(postID int, categories []int) error {
	for _, categorie := range categories {
		stmt, err := r.db.Prepare("INSERT INTO posts_categories (post_id, categorie_id) VALUES (?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(&postID, &categorie)
		if err != nil {
			return err
		}
	}

	return nil
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
