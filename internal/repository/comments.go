package repository

import (
	"database/sql"

	"github.com/alseiitov/real-time-forum/internal/model"
)

type CommentsRepo struct {
	db *sql.DB
}

func NewCommentsRepo(db *sql.DB) *CommentsRepo {
	return &CommentsRepo{db: db}
}

func (r *CommentsRepo) Create(comment model.Comment) (int, error) {
	stmt, err := r.db.Prepare("INSERT INTO comments (status, user_id, post_id, data, image, date) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(&comment.Status, &comment.UserID, &comment.PostID, &comment.Data, &comment.Image, &comment.Date)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	return int(id), err
}

func (r *CommentsRepo) Delete(userID, commentID int) error {
	res, err := r.db.Exec("DELETE FROM comments WHERE (id=$1) and (user_id=$2 OR EXISTS (SELECT * FROM users WHERE id=$2 AND role=$3))", commentID, userID, model.Roles.Admin)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if n == 0 {
		return ErrDeletingComment
	}

	return err
}

func (r *CommentsRepo) GetCommentsByPostID(postID int, limit int, offset int) ([]model.Comment, error) {
	var comments []model.Comment

	rows, err := r.db.Query("SELECT id, user_id, post_id, data, image, date FROM comments WHERE post_id = $1 LIMIT $2 OFFSET $3", postID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment model.Comment
		err = rows.Scan(&comment.ID, &comment.UserID, &comment.PostID, &comment.Data, &comment.Image, &comment.Date)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, rows.Err()
}
