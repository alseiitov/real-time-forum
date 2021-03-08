package handler

import (
	"net/http"

	"github.com/alseiitov/gorouter"
)

func (h *Handler) getRequestsForModerator(ctx *gorouter.Context) {
	requesters, err := h.adminsService.GetModeratorRequesters()
	if err != nil {
		ctx.WriteError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.WriteJSON(http.StatusOK, &requesters)
}

func (h *Handler) AcceptRequestForModerator(ctx *gorouter.Context) {
	requesterID, err := ctx.GetIntParam("requester_id")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	err = h.adminsService.AcceptRequestForModerator(requesterID)
	if err != nil {
		ctx.WriteError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.WriteHeader(http.StatusOK)
}

func (h *Handler) DeclineRequestForModerator(ctx *gorouter.Context) {
	requesterID, err := ctx.GetIntParam("requester_id")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	err = h.adminsService.DeclineRequestForModerator(requesterID)
	if err != nil {
		ctx.WriteError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.WriteHeader(http.StatusOK)
}
