package http

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/alseiitov/gorouter"
	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/pkg/auth"
)

func (h *Handler) cors(next gorouter.Handler) gorouter.Handler {
	return func(ctx *gorouter.Context) {
		(ctx.ResponseWriter).Header().Set("Access-Control-Allow-Origin", "*")
		(ctx.ResponseWriter).Header().Set("Access-Control-Allow-Methods", "*")
		(ctx.ResponseWriter).Header().Set("Access-Control-Allow-Headers", "*")

		if ctx.Request.Method == http.MethodOptions {
			ctx.WriteHeader(http.StatusOK)
			return
		}

		next(ctx)
	}
}

func (h *Handler) identify(minRole int, next gorouter.Handler) gorouter.Handler {
	return func(ctx *gorouter.Context) {
		token := ctx.Request.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
		if minRole > model.Roles.Guest || token != "" {
			sub, role, statusCode, err := h.identifyByToken(token, minRole)
			if err != nil {
				ctx.WriteError(statusCode, err.Error())
				return
			}
			ctx.SetParam("sub", strconv.Itoa(sub))
			ctx.SetParam("role", strconv.Itoa(role))
		}
		next(ctx)
	}
}

func (h *Handler) identifyByToken(token string, minRole int) (sub int, role int, statusCode int, err error) {
	if token == "" {
		if minRole > model.Roles.Guest {
			statusCode = http.StatusUnauthorized
			err = errors.New("401 Unauthorized")
			return
		}
		role = model.Roles.Guest
		return
	} else {
		sub, role, err = h.tokenManager.Parse(token)

		if err != nil {
			if err == auth.ErrExpiredToken {
				statusCode = http.StatusUnauthorized
			} else {
				statusCode = http.StatusBadRequest
			}
			return
		}

		if role < minRole {
			statusCode = http.StatusForbidden
			err = errors.New("HTTP 403 Forbidden")
			return
		}
	}
	return
}
