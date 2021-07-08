package http

import (
	"log"
	"net/http"

	"github.com/alseiitov/gorouter"
	_ "github.com/alseiitov/real-time-forum/docs"
	"github.com/alseiitov/real-time-forum/internal/handler/ws"
	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/service"
	"github.com/alseiitov/real-time-forum/pkg/auth"
	httpSwagger "github.com/swaggo/http-swagger"
)

type route struct {
	Path    string
	Method  string
	MinRole int
	Handler gorouter.Handler
}

type Handler struct {
	Router               *gorouter.Router
	wsHandler            *ws.Handler
	usersService         service.Users
	moderatorsService    service.Moderators
	adminsService        service.Admins
	categoriesService    service.Categories
	postsService         service.Posts
	commentsService      service.Comments
	notificationsService service.Notifications
	chatsService         service.Chats
	tokenManager         auth.TokenManager
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager, wsHandler *ws.Handler) *Handler {
	return &Handler{
		Router:               gorouter.NewRouter(),
		wsHandler:            wsHandler,
		usersService:         services.Users,
		moderatorsService:    services.Moderators,
		adminsService:        services.Admins,
		categoriesService:    services.Categories,
		postsService:         services.Posts,
		commentsService:      services.Comments,
		notificationsService: services.Notifications,
		chatsService:         services.Chats,
		tokenManager:         tokenManager,
	}
}

func (h *Handler) Init() {
	images := http.FileServer(http.Dir("./database/images"))

	go h.wsHandler.LogConns()

	routes := []route{
		// User handlers
		{
			Path:    "/api/users/sign-up",
			Method:  "POST",
			MinRole: model.Roles.Guest,
			Handler: h.usersSignUp,
		},
		{
			Path:    "/api/users/sign-in",
			Method:  "POST",
			MinRole: model.Roles.Guest,
			Handler: h.usersSignIn,
		},
		{
			Path:    "/api/users/:user_id",
			Method:  "GET",
			MinRole: model.Roles.User,
			Handler: h.getUser,
		},
		{
			Path:    "/api/users/:user_id/posts",
			Method:  "GET",
			MinRole: model.Roles.User,
			Handler: h.getUsersPosts,
		},
		{
			Path:    "/api/users/:user_id/rated-posts",
			Method:  "GET",
			MinRole: model.Roles.User,
			Handler: h.getUsersRatedPosts,
		},
		{
			Path:    "/api/auth/refresh",
			Method:  "POST",
			MinRole: model.Roles.Guest,
			Handler: h.usersRefreshTokens,
		},
		{
			Path:    "/api/moderators/requests",
			Method:  "POST",
			MinRole: model.Roles.User,
			Handler: h.requestModerator,
		},
		{
			Path:    "/api/notifications",
			Method:  "GET",
			MinRole: model.Roles.User,
			Handler: h.getNotifications,
		},

		// Post handlers
		{
			Path:    "/api/posts/:post_id",
			Method:  "GET",
			MinRole: model.Roles.Guest,
			Handler: h.getPost,
		},
		{
			Path:    "/api/posts",
			Method:  "POST",
			MinRole: model.Roles.User,
			Handler: h.createPost,
		},
		{
			Path:    "/api/posts/:post_id",
			Method:  "DELETE",
			MinRole: model.Roles.User,
			Handler: h.deletePost,
		},
		{
			Path:    "/api/posts/:post_id/likes",
			Method:  "POST",
			MinRole: model.Roles.User,
			Handler: h.likePost,
		},

		// Categories handlers
		{
			Path:    "/api/categories",
			Method:  "GET",
			MinRole: model.Roles.Guest,
			Handler: h.getAllCategories,
		},
		{
			Path:    "/api/categories/:category_id/:page",
			Method:  "GET",
			MinRole: model.Roles.Guest,
			Handler: h.getCategoryPage,
		},

		// Comments Handlers
		{
			Path:    "/api/posts/:post_id/comments/:page",
			Method:  "GET",
			MinRole: model.Roles.Guest,
			Handler: h.getCommentsOfPost,
		},
		{
			Path:    "/api/posts/:post_id/comments",
			Method:  "POST",
			MinRole: model.Roles.User,
			Handler: h.createComment,
		},
		{
			Path:    "/api/comments/:comment_id/likes",
			Method:  "POST",
			MinRole: model.Roles.User,
			Handler: h.likeComment,
		},
		{
			Path:    "/api/comments/:comment_id",
			Method:  "DELETE",
			MinRole: model.Roles.User,
			Handler: h.deleteComment,
		},

		// Admins Handlers
		{
			Path:    "/api/moderators/requests",
			Method:  "GET",
			MinRole: model.Roles.Admin,
			Handler: h.getRequestsForModerator,
		},
		{
			Path:    "/api/moderators/requests/:request_id",
			Method:  "POST",
			MinRole: model.Roles.Admin,
			Handler: h.RequestForModeratorAction,
		},

		// Chat Handlers
		{
			Path:    "/ws",
			Method:  "GET",
			MinRole: model.Roles.Guest,
			Handler: h.wsHandler.ServeWS,
		},

		// Swagger handler
		{
			Path:    "/swagger/*",
			Method:  "GET",
			MinRole: model.Roles.Guest,
			Handler: gorouter.WrapHandler(httpSwagger.Handler(httpSwagger.URL("http://localhost:8081/swagger/doc.json"))),
		},

		// Images filserver
		{
			Path:    "/images/*",
			Method:  "GET",
			MinRole: model.Roles.Guest,
			Handler: gorouter.WrapHandler(http.StripPrefix("/images/", images)),
		},
	}

	for _, route := range routes {
		switch route.Method {
		case "GET":
			h.Router.GET(route.Path, h.cors(h.identify(route.MinRole, route.Handler)))
		case "POST":
			h.Router.POST(route.Path, h.cors(h.identify(route.MinRole, route.Handler)))
		case "DELETE":
			h.Router.DELETE(route.Path, h.cors(h.identify(route.MinRole, route.Handler)))
		default:
			log.Fatalf("error: invalid method \"%v\" for route \"%v\"", route.Method, route.Path)
		}
	}

	go h.wsHandler.RunEventsPump()
}
