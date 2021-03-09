package app

import (
	"log"
	"os"

	"github.com/alseiitov/real-time-forum/pkg/auth"
	"github.com/alseiitov/real-time-forum/pkg/hash"

	"github.com/alseiitov/gorouter"
	"github.com/alseiitov/real-time-forum/internal/config"
	"github.com/alseiitov/real-time-forum/internal/handler"
	"github.com/alseiitov/real-time-forum/internal/repository"
	"github.com/alseiitov/real-time-forum/internal/server"
	"github.com/alseiitov/real-time-forum/internal/service"
	"github.com/alseiitov/real-time-forum/pkg/database"
)

func Run(configPath *string) {
	// Get forum config
	config, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatalln(err)
	}

	// Prepare database
	dbDriver := config.DBDriver()
	dbPath := config.DBPath()
	dbFileName := config.DBFileName()
	dbSchemesDir := config.DBSchemesDir()

	db, err := database.ConnectDB(dbDriver, dbPath, dbFileName, dbSchemesDir)
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
	accessTokenTTL, refreshTokenTTL, err := config.TokenTTLs()
	if err != nil {
		log.Fatalln(err)
	}

	tokenManager, err := auth.NewManager(jwtSigningKey, accessTokenTTL, refreshTokenTTL)
	if err != nil {
		log.Fatalln(err)
	}

	// Prepare services
	imagesDir := config.ImagesDir()
	defaultMaleAvatar, defaultFemaleAvatar := config.DefaultAvatars()
	postsForPage := config.PostsForPage()
	commentsForPage := config.CommentsForPage()
	postsPreModerationIsEnabled := config.PostsPreModerationIsEnabled()
	commentsPreModerationIsEnabled := config.CommentsPreModerationIsEnabled()

	services := service.NewServices(service.ServicesDeps{
		Repos:                          repos,
		Hasher:                         hasher,
		TokenManager:                   tokenManager,
		AccessTokenTTL:                 accessTokenTTL,
		RefreshTokenTTL:                refreshTokenTTL,
		ImagesDir:                      imagesDir,
		DefaultMaleAvatar:              defaultMaleAvatar,
		DefaultFemaleAvatar:            defaultFemaleAvatar,
		PostsForPage:                   postsForPage,
		CommentsForPage:                commentsForPage,
		PostsPreModerationIsEnabled:    postsPreModerationIsEnabled,
		CommentsPreModerationIsEnabled: commentsPreModerationIsEnabled,
	})

	// Prepare router
	router := gorouter.NewRouter()

	handler := handler.NewHandler(services.Users, services.Moderators, services.Admins, services.Categories, services.Posts, services.Comments, tokenManager)
	handler.Init(router)

	// Run server
	server := server.NewServer(config, router)
	log.Fatalln(server.Run())
}
