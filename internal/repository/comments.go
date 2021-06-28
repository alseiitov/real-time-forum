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
	stmt, err := r.db.Prepare(`
		INSERT INTO 
			comments (status, user_id, post_id, data, image, date) 
		VALUES 
			(?, ?, ?, ?, ?, ?)`,
	)

	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		&comment.Status,
		&comment.Author.ID,
		&comment.PostID,
		&comment.Data,
		&comment.Image,
		&comment.Date,
	)

	if err != nil {
		if isForeignKeyConstraintError(err) {
			return 0, ErrForeignKeyConstraint
		}
		return 0, err
	}

	id, err := res.LastInsertId()
	return int(id), err
}

func (r *CommentsRepo) GetByID(commentID int) (model.Comment, error) {
	var comment model.Comment

	row := r.db.QueryRow(`
		SELECT
			id, status, user_id, post_id, data, image, date
		FROM
			comments
		WHERE
			id = $1
		`,
		commentID,
	)

	err := row.Scan(
		&comment.ID,
		&comment.Status,
		&comment.Author.ID,
		&comment.PostID,
		&comment.Data,
		&comment.Image,
		&comment.Date,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return comment, ErrNoRows
		}
		return comment, err
	}

	return comment, nil
}

func (r *CommentsRepo) Delete(userID, commentID int) error {
	res, err := r.db.Exec(`DELETE FROM comments WHERE id=$1 AND user_id=$2`, commentID, userID)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if n == 0 {
		return ErrNoRows
	}

	return err
}

func (r *CommentsRepo) GetCommentsByPostID(postID int, userID int, limit int, offset int) ([]model.Comment, error) {
	var comments []model.Comment
	var postExists bool

	r.db.QueryRow(`SELECT EXISTS (SELECT id FROM posts WHERE id = $1)`, postID).Scan(&postExists)
	if !postExists {
		return nil, ErrNoRows
	}

	rows, err := r.db.Query(`
		SELECT
			comments.id, comments.post_id, 
			comments.user_id AS author_id, u.first_name AS author_first_name, u.last_name AS author_last_name, 
			comments.data, comments.image, comments.date, 
			IFNULL(cr.type, 0) AS user_rate, 
			COUNT(DISTINCT cl.id) - COUNT(DISTINCT cd.id) AS rating
		
		FROM 
			comments 
		LEFT JOIN users u 
		ON u.id == comments.user_id
			
		LEFT JOIN comments_likes cr 
		ON cr.comment_id = comments.id 
		AND cr.user_id = $1 
		
		LEFT JOIN comments_likes cl 
		ON cl.comment_id = comments.id 
		AND cl.type = 1 

		LEFT JOIN comments_likes cd 
		ON cd.comment_id = comments.id 
		AND cd.type = 2 
		
		WHERE comments.post_id = $2 
		GROUP BY comments.id 
		ORDER BY comments.id DESC
		LIMIT $3 OFFSET $4
		`,
		userID, postID, limit, offset,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment model.Comment

		err = rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.Author.ID,
			&comment.Author.FirstName,
			&comment.Author.LastName,
			&comment.Data,
			&comment.Image,
			&comment.Date,
			&comment.UserRate,
			&comment.Rating,
		)

		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	return comments, rows.Err()
}

func (r *CommentsRepo) LikeComment(like model.CommentLike) (bool, error) {
	var likeCreated bool // bool for checking that like/dislike created (not unliked/undisliked)

	tx, err := r.db.Begin()
	if err != nil {
		return likeCreated, err
	}

	// Get old like to comapare with new one
	var oldLike model.CommentLike
	row := tx.QueryRow(`
		SELECT 
			id, comment_id, user_id, type 
		FROM 
			comments_likes 
		WHERE 
			comment_id = $1 
		AND 
			user_id = $2
		`,
		like.CommentID, like.UserID,
	)

	err = row.Scan(
		&oldLike.ID,
		&oldLike.CommentID,
		&oldLike.UserID,
		&oldLike.LikeType,
	)

	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return likeCreated, err
	}

	// Delete old like if user already like this post
	if err == nil {
		_, err := tx.Exec(`DELETE FROM comments_likes WHERE id = $1`, oldLike.ID)
		if err != nil {
			tx.Rollback()
			return likeCreated, err
		}
	}

	// Create new like if user didn't like this comment or if type of new like and old like are not the same
	if err == sql.ErrNoRows || like.LikeType != oldLike.LikeType {
		_, err = tx.Exec(`
			INSERT INTO 
				comments_likes (comment_id, user_id, type) 
			VALUES 
				($1, $2, $3)
			`,
			like.CommentID, like.UserID, like.LikeType,
		)

		if err != nil {
			tx.Rollback()
			if isForeignKeyConstraintError(err) {
				return likeCreated, ErrForeignKeyConstraint
			}
			return likeCreated, err
		}

		likeCreated = true
	}

	return likeCreated, tx.Commit()
}
