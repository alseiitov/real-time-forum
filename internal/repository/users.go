package repository

import (
	"database/sql"

	"github.com/alseiitov/real-time-forum/internal/model"
)

type UsersRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UsersRepo {
	return &UsersRepo{db: db}
}

func (r *UsersRepo) Create(user model.User) error {
	stmt, err := r.db.Prepare(`
		INSERT INTO 
			users 
				(username, first_name, last_name, age, gender, email, password, role, avatar, registered) 
			VALUES 
				(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
	)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		user.Username,
		user.FirstName,
		user.LastName,
		user.Age,
		user.Gender,
		user.Email,
		user.Password,
		user.Role,
		user.Avatar,
		user.Registered,
	)

	if isAlreadyExistError(err) {
		return ErrAlreadyExist
	}

	return err
}

func (r *UsersRepo) GetByCredentials(usernameOrEmail, password string) (model.User, error) {
	var user model.User

	row := r.db.QueryRow(`
		SELECT 
			id, role 
		FROM 
			users 
		WHERE 
			(username = $1 OR email = $1) 
		AND 
			(password = $2)
		`,
		usernameOrEmail, password,
	)

	err := row.Scan(&user.ID, &user.Role)

	if isNoRowsError(err) {
		return user, ErrNoRows
	}

	return user, err
}

func (r *UsersRepo) GetByID(userID int) (model.User, error) {
	var user model.User

	row := r.db.QueryRow(`
		SELECT 
			id, username, first_name, last_name, age, gender, role, avatar, registered 
		FROM 
			users 
		WHERE 
			id = $1
		`,
		userID,
	)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Age,
		&user.Gender,
		&user.Role,
		&user.Avatar,
		&user.Registered,
	)

	if isNoRowsError(err) {
		return user, ErrNoRows
	}

	return user, err
}

func (r *UsersRepo) GetUsersPosts(userID int) ([]model.Post, error) {
	var userExists bool

	r.db.QueryRow(`SELECT EXISTS (SELECT id FROM users WHERE id = $1)`, userID).Scan(&userExists)
	if !userExists {
		return nil, ErrNoRows
	}

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
			posts.user_id = $1
		ORDER BY posts.id DESC
		`,
		userID,
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

		posts = append(posts, post)
	}

	return posts, nil
}

func (r *UsersRepo) GetUsersRatedPosts(userID int) ([]model.Post, error) {
	var userExists bool

	r.db.QueryRow(`SELECT EXISTS (SELECT id FROM users WHERE id = $1)`, userID).Scan(&userExists)
	if !userExists {
		return nil, ErrNoRows
	}

	var posts []model.Post

	rows, err := r.db.Query(`
		SELECT
			posts.id,
			posts.user_id AS author_id,
			u.first_name AS author_first_name,
			u.last_name AS author_last_name,
			posts.title,
			posts.date,
            pl.type AS user_rate
		FROM
			posts
			LEFT JOIN users u ON posts.user_id = u.id
            LEFT JOIN posts_likes pl ON posts.id = pl.post_id AND pl.user_id = $1  
		WHERE
			posts.id IN (
            	SELECT 
              		post_id 
              	FROM 
              		posts_likes 
              	WHERE 
              		user_id = $1
            )
		ORDER BY posts.id DESC
		`,
		userID,
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
			&post.UserRate,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (r *UsersRepo) CreateModeratorRequest(userID int) error {
	_, err := r.db.Exec(`INSERT INTO moderator_requests (user_id) VALUES ($1)`, userID)

	if isAlreadyExistError(err) {
		return ErrAlreadyExist
	}

	return err
}
