package handler

import (
	"github.com/alseiitov/gorouter"
	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/service"
	"github.com/alseiitov/real-time-forum/pkg/auth"
)

type Handler struct {
	usersService service.Users
	postsService service.Posts
	tokenManager auth.TokenManager
}

func NewHandler(usersService service.Users, postsService service.Posts, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		usersService: usersService,
		postsService: postsService,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init(r *gorouter.Router) {
	r.POST("/api/users/sign-up",
		h.cors(h.identify(model.Roles.Guest, h.usersSignUp)))

	r.POST("/api/users/sign-in",
		h.cors(h.identify(model.Roles.Guest, h.usersSignIn)))

	r.GET("/api/users/:id",
		h.cors(h.identify(model.Roles.User, h.getUser)))

	// r.PATCH("/api/users/:id",
	// 	h.cors(h.updateUser))

	r.POST("/api/auth/refresh",
		h.cors(h.usersRefreshToken))

	r.GET("/api/posts",
		h.cors(h.getAllPosts))

	r.POST("/api/posts",
		h.cors(h.identify(model.Roles.User, h.createPost)))

	r.GET("/api/posts/:id",
		h.cors(h.getPost))

	// r.PATCH("/api/posts/:id", h.cors(h.updatePost))
}
