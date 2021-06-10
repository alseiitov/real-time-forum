package handler

import (
	"log"
	"net/http"

	"github.com/alseiitov/gorouter"
	_ "github.com/alseiitov/real-time-forum/docs"
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
	eventsChan           chan *model.WSEvent
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

func NewHandler(services *service.Services, tokenManager auth.TokenManager, eventsChan chan *model.WSEvent) *Handler {
	return &Handler{
		Router:               gorouter.NewRouter(),
		eventsChan:           eventsChan,
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

	go logConns()

	routes := []route{
		// User handlers
		route{
			Path:    "/api/users/sign-up",
			Method:  "POST",
			MinRole: model.Roles.Guest,
			Handler: h.usersSignUp,
		},
		route{
			Path:    "/api/users/sign-in",
			Method:  "POST",
			MinRole: model.Roles.Guest,
			Handler: h.usersSignIn,
		},
		route{
			Path:    "/api/users/:user_id",
			Method:  "GET",
			MinRole: model.Roles.User,
			Handler: h.getUser,
		},
		route{
			Path:    "/api/auth/refresh",
			Method:  "POST",
			MinRole: model.Roles.Guest,
			Handler: h.usersRefreshTokens,
		},
		route{
			Path:    "/api/moderators/requests",
			Method:  "POST",
			MinRole: model.Roles.User,
			Handler: h.requestModerator,
		},
		route{
			Path:    "/api/notifications",
			Method:  "GET",
			MinRole: model.Roles.User,
			Handler: h.getNotifications,
		},

		// Post handlers
		route{
			Path:    "/api/posts/:post_id",
			Method:  "GET",
			MinRole: model.Roles.Guest,
			Handler: h.getPost,
		},
		route{
			Path:    "/api/posts",
			Method:  "POST",
			MinRole: model.Roles.User,
			Handler: h.createPost,
		},
		route{
			Path:    "/api/posts/:post_id",
			Method:  "DELETE",
			MinRole: model.Roles.User,
			Handler: h.deletePost,
		},
		route{
			Path:    "/api/posts/:post_id/likes",
			Method:  "POST",
			MinRole: model.Roles.User,
			Handler: h.likePost,
		},

		// Categories handlers
		route{
			Path:    "/api/categories",
			Method:  "GET",
			MinRole: model.Roles.Guest,
			Handler: h.getAllCategories,
		},
		route{
			Path:    "/api/categories/:category_id/:page",
			Method:  "GET",
			MinRole: model.Roles.Guest,
			Handler: h.getCategoryPage,
		},

		// Comments Handlers
		route{
			Path:    "/api/posts/:post_id/comments/:page",
			Method:  "GET",
			MinRole: model.Roles.Guest,
			Handler: h.getCommentsOfPost,
		},
		route{
			Path:    "/api/posts/:post_id/comments",
			Method:  "POST",
			MinRole: model.Roles.User,
			Handler: h.createComment,
		},
		route{
			Path:    "/api/comments/:comment_id/likes",
			Method:  "POST",
			MinRole: model.Roles.User,
			Handler: h.likeComment,
		},
		route{
			Path:    "/api/comments/:comment_id",
			Method:  "DELETE",
			MinRole: model.Roles.User,
			Handler: h.deleteComment,
		},

		// Admins Handlers
		route{
			Path:    "/api/moderators/requests",
			Method:  "GET",
			MinRole: model.Roles.Admin,
			Handler: h.getRequestsForModerator,
		},
		route{
			Path:    "/api/moderators/requests/:request_id",
			Method:  "POST",
			MinRole: model.Roles.Admin,
			Handler: h.RequestForModeratorAction,
		},

		// Chat Handlers
		route{
			Path:    "/ws",
			Method:  "GET",
			MinRole: model.Roles.Guest,
			Handler: h.handleWebSocket,
		},

		// Swagger handler
		route{
			Path:    "/swagger/*",
			Method:  "GET",
			MinRole: model.Roles.Guest,
			Handler: gorouter.WrapHandler(httpSwagger.Handler(httpSwagger.URL("http://localhost:8081/swagger/doc.json"))),
		},

		// Images filserver
		route{
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

	go h.runEventsPump()
}
