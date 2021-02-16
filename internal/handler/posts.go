package handler

import (
	"net/http"
	"strconv"

	"github.com/alseiitov/gorouter"
	"github.com/alseiitov/real-time-forum/internal/service"
	"github.com/alseiitov/validator"
)

func (h *Handler) getAllPosts(ctx *gorouter.Context) {

}

func (h *Handler) getPost(ctx *gorouter.Context) {

}

type createPostInput struct {
	Title      string `json:"title" validator:"required,min=2, max=64"`
	Data       string `json:"data" validator:"required,min=2, max=512"`
	Categories []int  `json:"categories" validator:"required,min=0"`
}

type createPostResponse struct {
	PostID int `json:"postID"`
}

func (h *Handler) createPost(ctx *gorouter.Context) {
	var input createPostInput
	sub, _ := ctx.GetParam("sub")
	userID, err := strconv.Atoi(sub)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	err = ctx.ReadBody(&input)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	err = validator.Validate(input)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	newPostID, err := h.postsService.Create(service.CreatePostInput{
		UserID:     userID,
		Title:      input.Title,
		Data:       input.Data,
		Categories: input.Categories,
	})
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	resp := createPostResponse{PostID: newPostID}
	ctx.WriteJSON(http.StatusCreated, resp)
}

type createCommentInput struct {
	Data string `json:"data" validator:"required,min=2,max=128"`
}

type createCommentResponse struct {
	CommentID int `json:"commentID"`
}

func (h *Handler) createComment(ctx *gorouter.Context) {
	var input createCommentInput

	sub, _ := ctx.GetParam("sub")
	userID, err := strconv.Atoi(sub)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	id, _ := ctx.GetParam("post_id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	err = ctx.ReadBody(&input)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	err = validator.Validate(input)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	newCommentID, err := h.postsService.CreateComment(service.CreateCommentInput{
		UserID: userID,
		PostID: postID,
		Data:   input.Data,
	})
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	resp := createCommentResponse{CommentID: newCommentID}
	ctx.WriteJSON(http.StatusCreated, resp)
}
