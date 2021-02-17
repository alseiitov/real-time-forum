package service

import (
	"time"

	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/repository"
)

const (
	withComments    = true
	withoutComments = false
)

type PostsService struct {
	repo repository.Posts
}

func NewPostsService(repo repository.Posts) *PostsService {
	return &PostsService{
		repo: repo,
	}
}

func (s *PostsService) Create(input CreatePostInput) (int, error) {
	post := model.Post{
		UserID: input.UserID,
		Title:  input.Title,
		Data:   input.Data,
		Date:   time.Now(),
		//TODO: add image upload
		Image:      "",
		Categories: input.Categories,
	}

	id, err := s.repo.Create(post)
	return id, err
}

func (s *PostsService) GetByID(role int, postID int) (model.Post, error) {
	var post model.Post
	var err error

	if role > model.Roles.Guest {
		post, err = s.repo.GetByID(postID, withComments)

	} else {
		post, err = s.repo.GetByID(postID, withoutComments)
	}
	if err != nil {
		return post, err
	}

	return post, nil
}

func (s *PostsService) Delete(userID int, role int, postID int) error {

	return s.repo.Delete(postID)
}

func (s *PostsService) CreateComment(input CreateCommentInput) (int, error) {
	comment := model.Comment{
		UserID: input.UserID,
		PostID: input.PostID,
		Data:   input.Data,
		//TODO: add image upload
		// Image: "",
		Date: time.Now(),
	}

	id, err := s.repo.CreateComment(comment)
	return id, err
}
