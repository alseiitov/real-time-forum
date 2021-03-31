package handler

import (
	"net/http"

	"github.com/alseiitov/gorouter"
	"github.com/alseiitov/real-time-forum/internal/service"
	"github.com/alseiitov/validator"
)

func (h *Handler) getPost(ctx *gorouter.Context) {
	postID, err := ctx.GetIntParam("post_id")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	post, err := h.postsService.GetByID(postID)
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
	Categories []int  `json:"categories" validator:"required,min=1"`
}

type createPostResponse struct {
	PostID int `json:"postID"`
}

func (h *Handler) createPost(ctx *gorouter.Context) {
	var input createPostInput
	userID, err := ctx.GetIntParam("sub")
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
		Categories: input.Categories,
	})

	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	resp := createPostResponse{PostID: newPostID}
	ctx.WriteJSON(http.StatusCreated, &resp)
}

func (h *Handler) deletePost(ctx *gorouter.Context) {
	userID, err := ctx.GetIntParam("sub")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	postID, err := ctx.GetIntParam("post_id")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	if err = h.postsService.Delete(userID, postID); err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	ctx.WriteHeader(http.StatusNoContent)
}

type likePostInput struct {
	LikeType int `json:"likeType" validator:"required,min=1,max=2"`
}

func (h *Handler) likePost(ctx *gorouter.Context) {
	var input likePostInput

	userID, err := ctx.GetIntParam("sub")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	postID, err := ctx.GetIntParam("post_id")
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

	if err = h.postsService.LikePost(postID, userID, input.LikeType); err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	ctx.WriteHeader(http.StatusOK)
}
