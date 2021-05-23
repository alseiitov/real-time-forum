package handler

import (
	"net/http"

	"github.com/alseiitov/gorouter"
	_ "github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/service"
	"github.com/alseiitov/validator"
)

type createCommentInput struct {
	Data  string `json:"data" validator:"required,min=2,max=128"`
	Image string `json:"image"`
}

type createCommentResponse struct {
	CommentID int `json:"commentID"`
}

// @Summary Create comment
// @Security Auth
// @Tags comments
// @ModuleID createComment
// @Accept  json
// @Produce  json
// @Param post_id path int true "ID of post"
// @Param input body createCommentInput true "comment input"
// @Success 201 {object} createCommentResponse
// @Failure 400,401,403,404,500 {object} gorouter.Error
// @Failure default {object} gorouter.Error
// @Router /posts/{post_id}/comments [POST]
func (h *Handler) createComment(ctx *gorouter.Context) {
	var input createCommentInput

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

	newCommentID, err := h.commentsService.Create(service.CreateCommentInput{
		UserID: userID,
		PostID: postID,
		Data:   input.Data,
		Image:  input.Image,
	})

	if err != nil {
		if err == service.ErrPostDoesntExist {
			ctx.WriteError(http.StatusNotFound, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	resp := createCommentResponse{CommentID: newCommentID}
	ctx.WriteJSON(http.StatusCreated, resp)
}

// @Summary Delete comment
// @Security Auth
// @Tags comments
// @ModuleID deleteComment
// @Accept  json
// @Produce  json
// @Param comment_id path int true "ID of comment"
// @Success 204 {string} ok
// @Failure 400,401,403,404,500 {object} gorouter.Error
// @Failure default {object} gorouter.Error
// @Router /comments/{comment_id} [DELETE]
func (h *Handler) deleteComment(ctx *gorouter.Context) {
	userID, err := ctx.GetIntParam("sub")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	commentID, err := ctx.GetIntParam("comment_id")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	err = h.commentsService.Delete(userID, commentID)
	if err != nil {
		if err == service.ErrDeletingComment {
			ctx.WriteError(http.StatusBadRequest, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	ctx.WriteHeader(http.StatusNoContent)
}

// @Summary Get page with N comments of post
// @Security Auth
// @Tags comments
// @ModuleID getCommentsOfPost
// @Accept  json
// @Produce  json
// @Param post_id path int true "ID of post"
// @Param page path int true "page number"
// @Success 200 {object} []model.Comment
// @Failure 400,401,403,404,500 {object} gorouter.Error
// @Failure default {object} gorouter.Error
// @Router /posts/{post_id}/comments/{page} [GET]
func (h *Handler) getCommentsOfPost(ctx *gorouter.Context) {
	postID, err := ctx.GetIntParam("post_id")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	page, err := ctx.GetIntParam("page")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	comments, err := h.commentsService.GetCommentsByPostID(postID, page)
	if err != nil {
		if err == service.ErrPostDoesntExist {
			ctx.WriteError(http.StatusNotFound, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	ctx.WriteJSON(http.StatusOK, &comments)
}

type likeCommentInput struct {
	LikeType int `json:"likeType" validator:"required,min=1,max=2"`
}

// @Summary Like of dislike comment
// @Security Auth
// @Tags comments
// @ModuleID likeComment
// @Accept  json
// @Produce  json
// @Param comment_id path int true "ID of comment"
// @Param input body likeCommentInput true "like type: 1 - like, 2 - dislike"
// @Success 200 {string} string "ok"
// @Failure 400,401,403,404,500 {object} gorouter.Error
// @Failure default {object} gorouter.Error
// @Router /comments/{comment_id}/likes [POST]
func (h *Handler) likeComment(ctx *gorouter.Context) {
	var input likeCommentInput

	userID, err := ctx.GetIntParam("sub")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	commentID, err := ctx.GetIntParam("comment_id")
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

	err = h.commentsService.LikeComment(commentID, userID, input.LikeType)
	if err != nil {
		if err == service.ErrCommentDoesntExist {
			ctx.WriteError(http.StatusNotFound, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	ctx.WriteHeader(http.StatusOK)
}
