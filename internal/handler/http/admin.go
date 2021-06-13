package http

import (
	"net/http"

	"github.com/alseiitov/gorouter"
	_ "github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/service"
)

// @Summary Get list of requests for moderator role
// @Security Auth
// @Tags admins
// @ModuleID getRequestsForModerator
// @Accept  json
// @Produce  json
// @Success 200 {object} []model.ModeratorRequest
// @Failure 400,401,403,404,500 {object} gorouter.Error
// @Failure default {object} gorouter.Error
// @Router /moderators/requests [GET]
func (h *Handler) getRequestsForModerator(ctx *gorouter.Context) {
	requests, err := h.adminsService.GetModeratorRequests()
	if err != nil {
		ctx.WriteError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.WriteJSON(http.StatusOK, &requests)
}

type RequestForModeratorActionInput struct {
	Action  string `json:"action" example:"accept"`
	Message string `json:"message"`
}

// @Summary Accept or decline request for moderator role
// @Security Auth
// @Tags admins
// @ModuleID getRequestsForModerator
// @Accept  json
// @Produce  json
// @Param request_id path int true "ID of request"
// @Param input body RequestForModeratorActionInput true "action 'accept' to accept or 'decline' to decline""
// @Success 200 {string} string "ok"
// @Failure 400,401,403,404,500 {object} gorouter.Error
// @Failure default {object} gorouter.Error
// @Router /moderators/requests/{request_id} [POST]
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

	if err = ctx.ReadBody(&input); err != nil {
		ctx.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	switch input.Action {
	case "accept":
		err = h.adminsService.AcceptRequestForModerator(adminID, requestID)
	case "decline":
		err = h.adminsService.DeclineRequestForModerator(adminID, requestID, input.Message)
	default:
		ctx.WriteError(http.StatusBadRequest, "invalid action")
		return
	}

	if err != nil {
		if err == service.ErrModeratorRequestDoesntExist {
			ctx.WriteError(http.StatusNotFound, err.Error())
		} else {
			ctx.WriteError(http.StatusInternalServerError, err.Error())
		}
		return
	}

	ctx.WriteHeader(http.StatusOK)
}
