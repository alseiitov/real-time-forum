package handler

import (
	"net/http"
	"strconv"

	"github.com/alseiitov/real-time-forum/internal/model"

	"github.com/alseiitov/gorouter"
	"github.com/alseiitov/real-time-forum/internal/service"
	"github.com/alseiitov/validator"
)

func (h *Handler) getPost(ctx *gorouter.Context) {
	postIDParam, _ := ctx.GetParam("post_id")
	postID, err := strconv.Atoi(postIDParam)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	post, err := h.postsService.GetByID(postID)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	post.Comments, err = h.commentsService.GetCommentsByPostID(postID)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	ctx.WriteJSON(http.StatusOK, &post)
}

type createPostInput struct {
	Title      string `json:"title" validator:"required,min=2, max=64"`
	Data       string `json:"data" validator:"required,min=2, max=512"`
	Image      string `json:"image"`
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

	if err = ctx.ReadBody(&input); err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	if err = validator.Validate(input); err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	newPostID, err := h.postsService.Create(service.CreatePostInput{
		UserID:     userID,
		Title:      input.Title,
		Data:       input.Data,
		Image:      input.Image,
		Categories: model.CategorieFromInts(input.Categories),
	})

	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	resp := createPostResponse{PostID: newPostID}
	ctx.WriteJSON(http.StatusCreated, &resp)
}

func (h *Handler) deletePost(ctx *gorouter.Context) {
	sub, _ := ctx.GetParam("sub")
	userID, err := strconv.Atoi(sub)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	postIDParam, _ := ctx.GetParam("post_id")
	postID, err := strconv.Atoi(postIDParam)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	err = h.postsService.Delete(userID, postID)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}
}
