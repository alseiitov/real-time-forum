package handler

import (
	"net/http"

	"github.com/alseiitov/gorouter"
	_ "github.com/alseiitov/real-time-forum/docs"
	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/service"
	"github.com/alseiitov/real-time-forum/pkg/auth"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	Router               *gorouter.Router
	usersService         service.Users
	moderatorsService    service.Moderators
	adminsService        service.Admins
	categoriesService    service.Categories
	postsService         service.Posts
	commentsService      service.Comments
	notificationsService service.Notifications
	tokenManager         auth.TokenManager
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		Router:               gorouter.NewRouter(),
		usersService:         services.Users,
		moderatorsService:    services.Moderators,
		adminsService:        services.Admins,
		categoriesService:    services.Categories,
		postsService:         services.Posts,
		commentsService:      services.Comments,
		notificationsService: services.Notifications,
		tokenManager:         tokenManager,
	}
}

func (h *Handler) Init() {
	h.initUsersHandlers()
	h.initPostsHandlers()
	h.initCategoriesHandlers()
	h.initCommentsHandlers()
	h.initAdminsHandlers()

	r := h.Router
	r.GET("/swagger/*", wrapHandler(httpSwagger.Handler(httpSwagger.URL("http://localhost:8081/swagger/doc.json"))))
}

func wrapHandler(h http.Handler) gorouter.Handler {
	return func(ctx *gorouter.Context) {
		h.ServeHTTP(ctx.ResponseWriter, ctx.Request)
	}
}

func (h *Handler) initUsersHandlers() {
	r := h.Router

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
}

func (h *Handler) initPostsHandlers() {
	r := h.Router

	r.GET("/api/posts/:post_id",
		h.cors(h.identify(model.Roles.Guest, h.getPost)))

	r.POST("/api/posts",
		h.cors(h.identify(model.Roles.User, h.createPost)))

	r.DELETE("/api/posts/:post_id",
		h.cors(h.identify(model.Roles.User, h.deletePost)))

	r.POST("/api/posts/:post_id/likes",
		h.cors(h.identify(model.Roles.User, h.likePost)))
}

func (h *Handler) initCategoriesHandlers() {
	r := h.Router

	r.GET("/api/categories",
		h.cors(h.identify(model.Roles.Guest, h.getAllCategories)))

	r.GET("/api/categories/:category_id/:page",
		h.cors(h.identify(model.Roles.Guest, h.getCategoryPage)))
}

func (h *Handler) initCommentsHandlers() {
	r := h.Router

	r.GET("/api/posts/:post_id/comments/:page",
		h.cors(h.identify(model.Roles.Guest, h.getCommentsOfPost)))

	r.POST("/api/posts/:post_id/comments",
		h.cors(h.identify(model.Roles.User, h.createComment)))

	r.POST("/api/comments/:comment_id/likes",
		h.cors(h.identify(model.Roles.User, h.likeComment)))

	r.DELETE("/api/comments/:comment_id",
		h.cors(h.identify(model.Roles.User, h.deleteComment)))
}

func (h *Handler) initAdminsHandlers() {
	r := h.Router

	r.GET("/api/moderators/requests",
		h.cors(h.identify(model.Roles.Admin, h.getRequestsForModerator)))

	r.POST("/api/moderators/requests/:request_id",
		h.cors(h.identify(model.Roles.Admin, h.RequestForModeratorAction)))
}
