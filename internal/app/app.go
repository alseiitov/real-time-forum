package app

import (
	"log"
	"os"

	"github.com/alseiitov/real-time-forum/pkg/auth"
	"github.com/alseiitov/real-time-forum/pkg/hash"

	"github.com/alseiitov/real-time-forum/internal/config"
	"github.com/alseiitov/real-time-forum/internal/handler"
	"github.com/alseiitov/real-time-forum/internal/repository"
	"github.com/alseiitov/real-time-forum/internal/server"
	"github.com/alseiitov/real-time-forum/internal/service"
	"github.com/alseiitov/real-time-forum/pkg/database"
)

// @title real-time-forum API
// @version 1.0
// @description API Server for real-time-forum project

// @host localhost:8081
// @BasePath /api

// @securityDefinitions.apikey Auth
// @in header
// @name Authorization

func Run(configPath *string) {
	// Get forum config
	config, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatalln(err)
	}

	// Prepare database
	db, err := database.ConnectDB(
		config.Database.Driver,
		config.Database.Path,
		config.Database.FileName,
		config.Database.SchemesDir,
	)
	if err != nil {
		log.Fatalln(err)
	}

	// Run function that deletes expired sessions from database
	go repository.DeleteExpiredSessions(db)

	// Prepare repositories
	repos := repository.NewRepositories(db)

	// Prepare password hasher
	passwordSalt := os.Getenv("PASSWORD_SALT")
	hasher, err := hash.NewHasher(passwordSalt)
	if err != nil {
		log.Fatalln(err)
	}

	// Prepare JWT token manager
	jwtSigningKey := os.Getenv("JWT_SIGNING_KEY")
	accessTokenTTL := config.AccessTokenTTL()
	refreshTokenTTL := config.RefreshTokenTTL()
	if err != nil {
		log.Fatalln(err)
	}

	tokenManager, err := auth.NewManager(jwtSigningKey, accessTokenTTL, refreshTokenTTL)
	if err != nil {
		log.Fatalln(err)
	}

	// Prepare services
	services := service.NewServices(service.ServicesDeps{
		Repos:                          repos,
		Hasher:                         hasher,
		TokenManager:                   tokenManager,
		AccessTokenTTL:                 accessTokenTTL,
		RefreshTokenTTL:                refreshTokenTTL,
		MaxConnsForUser:                config.Websocket.MaxConnsForUser,
		MaxMessageSize:                 config.Websocket.MaxMessageSize,
		TokenWait:                      config.TokenWait(),
		WriteWait:                      config.WriteWait(),
		PongWait:                       config.PongWait(),
		PingPeriod:                     config.PingPeriod(),
		ImagesDir:                      config.Database.ImagesDir,
		DefaultMaleAvatar:              config.Forum.DefaultMaleAvatar,
		DefaultFemaleAvatar:            config.Forum.DefaultFemaleAvatar,
		PostsForPage:                   config.Forum.PostsForPage,
		CommentsForPage:                config.Forum.CommentsForPage,
		PostsPreModerationIsEnabled:    config.Forum.PostsPreModerationIsEnabled,
		CommentsPreModerationIsEnabled: config.Forum.CommentsPreModerationIsEnabled,
	})

	// Prepare handler
	handler := handler.NewHandler(services, tokenManager)
	handler.Init()

	// Run server
	server := server.NewServer(config, handler.Router)
	log.Fatalln(server.Run())
}
