package handler

import (
	"net/http"
	"strconv"

	"github.com/alseiitov/gorouter"
	"github.com/alseiitov/real-time-forum/internal/service"
	"github.com/alseiitov/validator"
)

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

	newCommentID, err := h.commentsService.Create(service.CreateCommentInput{
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

	err = h.commentsService.Delete(userID, commentID)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}
}
