package handler

import (
	"github.com/alseiitov/gorouter"
	"github.com/alseiitov/real-time-forum/internal/service"
)

type Handler struct {
	userService service.Users
}

func NewHandler(usersService service.Users) *Handler {
	return &Handler{
		userService: usersService,
	}
}

func (h *Handler) Init(router *gorouter.Router) {
	router.GET("/api/users/:id", h.usersSignUp)
}
