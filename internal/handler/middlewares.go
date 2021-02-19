package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/alseiitov/gorouter"
	"github.com/alseiitov/real-time-forum/internal/model"
)

func (h *Handler) cors(next gorouter.Handler) gorouter.Handler {
	return func(ctx *gorouter.Context) {
		(ctx.ResponseWriter).Header().Set("Access-Control-Allow-Origin", "*")
		(ctx.ResponseWriter).Header().Set("Access-Control-Allow-Methods", "POST, GET, PATCH, DELETE")
		(ctx.ResponseWriter).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")

		next(ctx)
	}
}

func (h *Handler) identify(minRole int, next gorouter.Handler) gorouter.Handler {
	return func(ctx *gorouter.Context) {
		token := ctx.Request.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
		if token == "" {
			if minRole > model.Roles.Guest {
				ctx.WriteError(http.StatusUnauthorized, "401 Unauthorized")
				return
			}

			ctx.SetParam("role", strconv.Itoa(model.Roles.Guest))
		} else {
			sub, role, err := h.tokenManager.Parse(token)

			if err != nil {
				ctx.WriteError(http.StatusBadRequest, err.Error())
				return
			}
			if role < minRole {
				ctx.WriteError(http.StatusForbidden, "HTTP 403 Forbidden")
				return
			}

			ctx.SetParam("sub", strconv.Itoa(sub))
			ctx.SetParam("role", strconv.Itoa(role))
		}
		next(ctx)
	}
}
