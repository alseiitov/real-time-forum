package handler

import (
	"net/http"

	"github.com/alseiitov/gorouter"
)

func (h *Handler) getRequestsForModerator(ctx *gorouter.Context) {
	requests, err := h.adminsService.GetModeratorRequests()
	if err != nil {
		ctx.WriteError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.WriteJSON(http.StatusOK, &requests)
}

func (h *Handler) AcceptRequestForModerator(ctx *gorouter.Context) {
	adminID, err := ctx.GetIntParam("sub")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	requestID, err := ctx.GetIntParam("request_id")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	err = h.adminsService.AcceptRequestForModerator(adminID, requestID)
	if err != nil {
		ctx.WriteError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.WriteHeader(http.StatusOK)
}

type DeclineRequestForModeratorInput struct {
	Message string `json:"message`
}

func (h *Handler) DeclineRequestForModerator(ctx *gorouter.Context) {
	var input DeclineRequestForModeratorInput

	adminID, err := ctx.GetIntParam("sub")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	requestID, err := ctx.GetIntParam("request_id")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	err = ctx.ReadBody(&input)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	err = h.adminsService.DeclineRequestForModerator(adminID, requestID, input.Message)
	if err != nil {
		ctx.WriteError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.WriteHeader(http.StatusOK)
}
