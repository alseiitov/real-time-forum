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

	ctx.WriteJSON(http.StatusOK, &post)
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
		Categories: model.CategorieFromInts(input.Categories...),
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

type createCommentInput struct {
	PostID int    `json:"postID" validator:"required"`
	Data   string `json:"data" validator:"required,min=2,max=128"`
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
		PostID: input.PostID,
		Data:   input.Data,
	})
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	resp := createCommentResponse{CommentID: newCommentID}
	ctx.WriteJSON(http.StatusCreated, resp)
}

func (h *Handler) deleteComment(ctx *gorouter.Context) {
	sub, _ := ctx.GetParam("sub")
	userID, err := strconv.Atoi(sub)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	commentIDParam, _ := ctx.GetParam("comment_id")
	commentID, err := strconv.Atoi(commentIDParam)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	err = h.postsService.DeleteComment(userID, commentID)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}
}

func (h *Handler) getCategories(ctx *gorouter.Context) {
	categories, err := h.postsService.GetCategories()
	if err != nil {
		ctx.WriteError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.WriteJSON(http.StatusOK, categories)
}

func (h *Handler) getPostsByCategory(ctx *gorouter.Context) {

}

type createCategoryInput struct {
	Name string `json:"name" validator:"required,min=2,max=128"`
}

type createCategoryResponse struct {
	CategoryID int `json:"categoryID"`
}

func (h *Handler) createCategory(ctx *gorouter.Context) {

}
