package handler

import (
	"github.com/alseiitov/gorouter"
	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/service"
	"github.com/alseiitov/real-time-forum/pkg/auth"
)

type Handler struct {
	usersService      service.Users
	categoriesService service.Categories
	postsService      service.Posts
	commentsService   service.Comments
	tokenManager      auth.TokenManager
}

func NewHandler(usersService service.Users, categoriesService service.Categories, postsService service.Posts, commentsService service.Comments, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		usersService:      usersService,
		categoriesService: categoriesService,
		postsService:      postsService,
		commentsService:   commentsService,
		tokenManager:      tokenManager,
	}
}

func (h *Handler) Init(r *gorouter.Router) {
	//
	//
	// Users handlers

	r.POST("/api/users/sign-up",
		h.cors(h.identify(model.Roles.Guest, h.usersSignUp)))

	r.POST("/api/users/sign-in",
		h.cors(h.identify(model.Roles.Guest, h.usersSignIn)))

	r.GET("/api/users/:user_id",
		h.cors(h.identify(model.Roles.Guest, h.getUser)))

	r.POST("/api/auth/refresh",
		h.cors(h.usersRefreshTokens))

	//
	//
	// Posts handlers

	r.GET("/api/posts/:post_id",
		h.cors(h.identify(model.Roles.Guest, h.getPost)))

	r.POST("/api/posts",
		h.cors(h.identify(model.Roles.User, h.createPost)))

	r.DELETE("/api/posts/:post_id",
		h.cors(h.identify(model.Roles.User, h.deletePost)))

	//
	//
	// Categories handlers

	r.GET("/api/categories",
		h.cors(h.identify(model.Roles.Guest, h.getAllCategories)))

	r.GET("/api/categories/:category_id/:page",
		h.cors(h.identify(model.Roles.Guest, h.getCategoryPage)))

	//
	//
	//Comments handlers

	r.POST("/api/comments",
		h.cors(h.identify(model.Roles.User, h.createComment)))

	r.DELETE("/api/comments/:comment_id",
		h.cors(h.identify(model.Roles.User, h.deleteComment)))

}
