package handler

import (
	"net/http"

	"github.com/alseiitov/gorouter"
	"github.com/alseiitov/real-time-forum/pkg/image"
)

func (h *Handler) getImage(ctx *gorouter.Context) {

}

type createImageResponse struct {
	Name string `json:"name"`
}

func (h *Handler) createImage(ctx *gorouter.Context) {
	file, stat, err := image.ParseFromRequest(ctx.Request)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	name, err := image.SaveImage(file, stat, "./database/images")
	if err != nil {
		ctx.WriteError(http.StatusInternalServerError, err.Error())
		return
	}

	resp := createImageResponse{Name: name}

	ctx.WriteJSON(http.StatusCreated, resp)
}
