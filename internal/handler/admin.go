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

type RequestForModeratorActionInput struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

func (h *Handler) RequestForModeratorAction(ctx *gorouter.Context) {
	var input RequestForModeratorActionInput

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

	switch input.Action {
	case "accept":
		err = h.adminsService.AcceptRequestForModerator(adminID, requestID)
	case "decline":
		err = h.adminsService.DeclineRequestForModerator(adminID, requestID, input.Message)
	default:
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	if err != nil {
		ctx.WriteError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.WriteHeader(http.StatusOK)
}
