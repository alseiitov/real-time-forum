package handler

import (
	"net/http"

	"github.com/alseiitov/gorouter"
	"github.com/alseiitov/real-time-forum/internal/service"
	"github.com/alseiitov/validator"
)

type createCommentInput struct {
	PostID int    `json:"postID" validator:"required,min=1"`
	Data   string `json:"data" validator:"required,min=2,max=128"`
	Image  string `json:"image"`
}

type createCommentResponse struct {
	CommentID int `json:"commentID"`
}

func (h *Handler) createComment(ctx *gorouter.Context) {
	var input createCommentInput

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

	newCommentID, err := h.commentsService.Create(service.CreateCommentInput{
		UserID: userID,
		PostID: input.PostID,
		Data:   input.Data,
		Image:  input.Image,
	})
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	resp := createCommentResponse{CommentID: newCommentID}
	ctx.WriteJSON(http.StatusCreated, resp)
}

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

	if err = h.commentsService.Delete(userID, commentID); err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	ctx.WriteHeader(http.StatusNoContent)
}

type getCommentsOfPostInput struct {
	PostID int `json:"postID" validator:"required,min=1"`
	Page   int `json:"page" validator:"required,min=1"`
}

func (h *Handler) getCommentsOfPost(ctx *gorouter.Context) {
	var input getCommentsOfPostInput

	if err := ctx.ReadBody(&input); err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	if err := validator.Validate(input); err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	comments, err := h.commentsService.GetCommentsByPostID(input.PostID, input.Page)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	ctx.WriteJSON(http.StatusOK, &comments)
}

type likeCommentInput struct {
	LikeType int `json:"likeType" validator:"required,min=1,max=2"`
}

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

	if err = h.commentsService.LikeComment(commentID, userID, input.LikeType); err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	ctx.WriteHeader(http.StatusOK)
}
