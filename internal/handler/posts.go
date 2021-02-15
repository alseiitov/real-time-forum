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
	Title      string   `json:"title" validator:"required,min=2, max=64"`
	Data       string   `json:"data" validator:"required,min=2, max=512"`
	Categories []string `json:"categories" validator:"required,min=5,max=64"`
}

func (h *Handler) createPost(ctx *gorouter.Context) {
	var input createPostInput
	sub, _ := ctx.GetParam("sub")
	userID, err := strconv.Atoi(sub)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	err := ctx.ReadBody(&input)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	err = validator.Validate(input)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	err = h.postsService.Create(service.CreatePostInput{
		UserID:     userID,
		Title:      input.Title,
		Data:       input.Data,
		Categories: input.Categories,
	})
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	ctx.ResponseWriter.WriteHeader(http.StatusCreated)
}
