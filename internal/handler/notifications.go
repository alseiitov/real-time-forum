package handler

import (
	"net/http"

	"github.com/alseiitov/gorouter"
)

func (h *Handler) getNotifications(ctx *gorouter.Context) {
	userID, err := ctx.GetIntParam("sub")
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	notifications, err := h.notificationsService.GetNotifications(userID)
	if err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
	}

	ctx.WriteJSON(http.StatusOK, notifications)
}
