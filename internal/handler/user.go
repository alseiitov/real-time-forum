package handler

import (
	"github.com/alseiitov/gorouter"
)

func (h *Handler) usersSignUp(ctx *gorouter.Context) {
	h.enableCORS(ctx)

	ctx.ResponseWriter.Write([]byte(ctx.Params[0]))
}
