package handler

import (
	"github.com/alseiitov/gorouter"
	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/service"
	"github.com/alseiitov/real-time-forum/pkg/auth"
)

type Handler struct {
	usersService         service.Users
	moderatorsService    service.Moderators
	adminsService        service.Admins
	categoriesService    service.Categories
	postsService         service.Posts
	commentsService      service.Comments
	notificationsService service.Notifications
	tokenManager         auth.TokenManager
}

func NewHandler(usersService service.Users, moderatorsService service.Moderators, adminsService service.Admins, categoriesService service.Categories, postsService service.Posts, commentsService service.Comments, notificationsService service.Notifications, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		usersService:         usersService,
		moderatorsService:    moderatorsService,
		adminsService:        adminsService,
		categoriesService:    categoriesService,
		postsService:         postsService,
		commentsService:      commentsService,
		notificationsService: notificationsService,
		tokenManager:         tokenManager,
	}
}

func (h *Handler) Init(r *gorouter.Router) {
	//
	// Users handlers
	//
	r.POST("/api/users/sign-up",
		h.cors(h.identify(model.Roles.Guest, h.usersSignUp)))

	r.POST("/api/users/sign-in",
		h.cors(h.identify(model.Roles.Guest, h.usersSignIn)))

	r.GET("/api/users/:user_id",
		h.cors(h.identify(model.Roles.Guest, h.getUser)))

	r.POST("/api/auth/refresh",
		h.cors(h.usersRefreshTokens))

	r.POST("/api/moderators/requests",
		h.cors(h.identify(model.Roles.User, h.requestModerator)))

	r.GET("/api/notifications",
		h.cors(h.identify(model.Roles.User, h.getNotifications)))

	//
	// Posts handlers
	//
	r.GET("/api/posts/:post_id",
		h.cors(h.identify(model.Roles.Guest, h.getPost)))

	r.POST("/api/posts",
		h.cors(h.identify(model.Roles.User, h.createPost)))

	r.DELETE("/api/posts/:post_id",
		h.cors(h.identify(model.Roles.User, h.deletePost)))

	r.POST("/api/posts/:post_id/likes",
		h.cors(h.identify(model.Roles.User, h.likePost)))

	//
	// Categories handlers
	//
	r.GET("/api/categories",
		h.cors(h.identify(model.Roles.Guest, h.getAllCategories)))

	r.GET("/api/categories/:category_id/:page",
		h.cors(h.identify(model.Roles.Guest, h.getCategoryPage)))

	//
	// Comments handlers
	//
	r.GET("/api/comments",
		h.cors(h.identify(model.Roles.Guest, h.getCommentsOfPost)))

	r.POST("/api/comments",
		h.cors(h.identify(model.Roles.User, h.createComment)))

	r.POST("/api/comments/:comment_id/likes",
		h.cors(h.identify(model.Roles.User, h.likeComment)))

	r.DELETE("/api/comments/:comment_id",
		h.cors(h.identify(model.Roles.User, h.deleteComment)))

	//
	// Admins handlers
	//
	r.GET("/api/moderators/requests",
		h.cors(h.identify(model.Roles.Admin, h.getRequestsForModerator)))

	r.POST("/api/moderators/requests/:request_id",
		h.cors(h.identify(model.Roles.Admin, h.RequestForModeratorAction)))

}
