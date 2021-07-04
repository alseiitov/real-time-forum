package repository

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

	stmt, err := tx.Prepare(`
		INSERT INTO 
			posts (status, user_id, title, data, date, image) 
		VALUES 
			(?, ?, ?, ?, ?, ?)`,
	)

	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		&post.Status,
		&post.Author.ID,
		&post.Title,
		&post.Data,
		&post.Date,
		&post.Image,
	)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, category := range post.Categories {
		stmt, err := tx.Prepare(`
			INSERT INTO 
				posts_categories (post_id, category_id) 
			VALUES 
				(?, ?)`,
		)

		if err != nil {
			tx.Rollback()
			return 0, err
		}

		_, err = stmt.Exec(&id, &category.ID)
		if err != nil {
			tx.Rollback()
			if isForeignKeyConstraintError(err) {
				return 0, ErrForeignKeyConstraint
			}
			return 0, err
		}
	}

	return int(id), tx.Commit()
}

func (r *PostsRepo) GetByID(postID int, userID int) (model.Post, error) {
	var post model.Post
	var postExists bool

	r.db.QueryRow(`SELECT EXISTS (SELECT id FROM posts WHERE id = $1)`, postID).Scan(&postExists)
	if !postExists {
		return post, ErrNoRows
	}

	row := r.db.QueryRow(`
		SELECT 
			posts.id, 
			posts.user_id AS author_id, 
			u.first_name AS author_first_name, 
			u.last_name AS author_last_name, 
			posts.title, 
			posts.data, 
			posts.date, 
			posts.image, 
			IFNULL(pr.type, 0) AS user_rate, 
			COUNT(DISTINCT pl.id) - COUNT(DISTINCT pd.id) AS rating 
		FROM 
			posts 
			LEFT JOIN users u 
			ON posts.user_id = u.id 
			
			LEFT JOIN posts_likes pr 
			ON pr.post_id = posts.id 
			AND pr.user_id = $1 

			LEFT JOIN posts_likes pl 
			ON pl.post_id = posts.id 
			AND pl.type = $2 
			
			LEFT JOIN posts_likes pd ON pd.post_id = posts.id 
			AND pd.type = $3 
		WHERE 
			posts.id = $4 
		`,
		userID, model.LikeTypes.Like, model.LikeTypes.Dislike, postID,
	)

	err := row.Scan(
		&post.ID,
		&post.Author.ID,
		&post.Author.FirstName,
		&post.Author.LastName,
		&post.Title,
		&post.Data,
		&post.Date,
		&post.Image,
		&post.UserRate,
		&post.Rating,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return post, ErrNoRows
		}
		return post, err
	}

	post.Categories, err = r.getPostCategories(postID)
	if err != nil {
		return post, err
	}

	return post, nil
}

func (r *PostsRepo) getPostCategories(postID int) ([]model.Category, error) {
	var categories []model.Category

	rows, err := r.db.Query(`
		SELECT 
			id, name 
		FROM 
			categories 
		WHERE 
			id IN (
				SELECT 
					category_id 
				FROM 
					posts_categories 
				WHERE 
					post_id = $1
			)
		`,
		postID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category model.Category
		err = rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, rows.Err()
}

func (r *PostsRepo) Delete(userID, postID int) error {
	res, err := r.db.Exec(`DELETE FROM posts WHERE id=$1 AND user_id=$2`, postID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNoRows
		}
		return err
	}

	n, err := res.RowsAffected()
	if n == 0 {
		return ErrNoRows
	}

	return err
}

func (r *PostsRepo) GetPostsByCategoryID(categoryID int, limit int, offset int) ([]model.Post, error) {
	var posts []model.Post

	rows, err := r.db.Query(`
		SELECT
			posts.id,
			posts.user_id AS author_id,
			u.first_name AS author_first_name,
			u.last_name AS author_last_name,
			posts.title,
			posts.date
		FROM
			posts
			LEFT JOIN users u ON posts.user_id = u.id
		WHERE
			posts.id IN (
				SELECT
					post_id
				FROM
					posts_categories
				WHERE
					category_id = $1
			)
		GROUP BY
			posts.id
		ORDER BY
			posts.id DESC
		LIMIT
			$2 OFFSET $3
		`,
		categoryID, limit, offset,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post model.Post
		err = rows.Scan(
			&post.ID,
			&post.Author.ID,
			&post.Author.FirstName,
			&post.Author.LastName,
			&post.Title,
			&post.Date,
		)
		if err != nil {
			return nil, err
		}

		post.Categories, err = r.getPostCategories(post.ID)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostsRepo) LikePost(like model.PostLike) (bool, error) {
	var likeCreated bool // bool for checking that like/dislike created (not unliked/undisliked)

	tx, err := r.db.Begin()
	if err != nil {
		return likeCreated, err
	}

	// Get old like to comapare with new one
	var oldLike model.PostLike
	row := tx.QueryRow(`
		SELECT 
			id, post_id, user_id, type 
		FROM 
			posts_likes 
		WHERE 
			post_id = $1 
		AND 
			user_id = $2
		`,
		like.PostID, like.UserID,
	)

	err = row.Scan(&oldLike.ID, &oldLike.PostID, &oldLike.UserID, &oldLike.LikeType)

	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return likeCreated, err
	}

	// Delete old like if user already like this post
	if err == nil {
		_, err := tx.Exec(`DELETE FROM posts_likes WHERE id = $1`, oldLike.ID)
		if err != nil {
			tx.Rollback()
			return likeCreated, err
		}
	}

	// Create new like if user didn't like this post or if type of new like and old like are not the same
	if err == sql.ErrNoRows || like.LikeType != oldLike.LikeType {
		_, err = tx.Exec(`
			INSERT INTO 
				posts_likes (post_id, user_id, type) 
			VALUES 
				($1, $2, $3)
			`,
			like.PostID, like.UserID, like.LikeType,
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
