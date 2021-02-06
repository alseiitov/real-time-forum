package handler

import (
	"github.com/alseiitov/gorouter"
)

func (h *Handler) enableCORS(ctx *gorouter.Context) {
	(ctx.ResponseWriter).Header().Set("Access-Control-Allow-Origin", "*")
	(ctx.ResponseWriter).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(ctx.ResponseWriter).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
}
