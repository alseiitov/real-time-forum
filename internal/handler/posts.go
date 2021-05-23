package handler

import (
	"net/http"

	"github.com/alseiitov/gorouter"
	_ "github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/service"
	"github.com/alseiitov/validator"
)

// @Summary Get post by ID
// @Security Auth
// @Tags posts
// @ModuleID getPost
// @Accept  json
// @Produce  json
// @Param post_id path int true "ID of post"
// @Success 200 {object} model.Post
// @Failure 400,404,500 {object} gorouter.Error
// @Failure default {object} gorouter.Error
// @Router /posts/{post_id} [GET]
func (h *Handler) getPost(ctx *gorouter.Context) {
	postID, err := ctx.GetIntParam("post_id")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	post, err := h.postsService.GetByID(postID)
	if err != nil {
		if err == service.ErrPostDoesntExist {
			ctx.WriteError(http.StatusNotFound, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
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

// @Summary Create post
// @Security Auth
// @Tags posts
// @ModuleID createPost
// @Accept  json
// @Produce  json
// @Param input body createPostInput true "post input data"
// @Success 201 {object} createPostResponse
// @Failure 400,401,404,500 {object} gorouter.Error
// @Failure default {object} gorouter.Error
// @Router /posts [POST]
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
		if err == service.ErrCategoryDoesntExist {
			ctx.WriteError(http.StatusNotFound, err.Error())
		} else if err == service.ErrTooManyCategories {
			ctx.WriteError(http.StatusBadRequest, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	resp := createPostResponse{PostID: newPostID}
	ctx.WriteJSON(http.StatusCreated, &resp)
}

// @Summary Delete post
// @Security Auth
// @Tags posts
// @ModuleID deletePost
// @Accept  json
// @Produce  json
// @Param post_id path int true "ID of post"
// @Success 204 {string} string "ok"
// @Failure 400,401,403,404,500 {object} gorouter.Error
// @Failure default {object} gorouter.Error
// @Router /posts/{post_id} [DELETE]
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
		if err == service.ErrDeletingPost {
			ctx.WriteError(http.StatusBadRequest, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	ctx.WriteHeader(http.StatusNoContent)
}

type likePostInput struct {
	LikeType int `json:"likeType" validator:"required,min=1,max=2"`
}

// @Summary Like or dislike post
// @Security Auth
// @Tags posts
// @ModuleID likePost
// @Accept  json
// @Produce  json
// @Param input body likePostInput true "like type: 1 - like, 2 - dislike"
// @Param post_id path int true "ID of post"
// @Success 200 {string} string "ok"
// @Failure 400,401,403,404,500 {object} gorouter.Error
// @Failure default {object} gorouter.Error
// @Router /posts/{post_id}/likes [POST]
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

	err = h.postsService.LikePost(postID, userID, input.LikeType)
	if err != nil {
		if err == service.ErrPostDoesntExist {
			ctx.WriteError(http.StatusNotFound, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	ctx.WriteHeader(http.StatusOK)
}
