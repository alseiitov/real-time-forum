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
	"github.com/alseiitov/real-time-forum/pkg/database/sqlite"
)

func Run(configPath *string) {
	config, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatalln(err)
	}

	db, err := sqlite.ConnectDB(config)
	if err != nil {
		log.Fatalln(err)
	}

	repos := repository.NewRepositories(db)

	passwordSalt := os.Getenv("PASSWORD_SALT")
	hasher, err := hash.NewHasher(passwordSalt)
	if err != nil {
		log.Fatalln(err)
	}

	jwtSigningKey := os.Getenv("JWT_SIGNING_KEY")
	accessTokenTTL, refreshTokenTTL, err := config.GetTokenTTLs()
	if err != nil {
		log.Fatalln(err)
	}

	tokenManager, err := auth.NewManager(jwtSigningKey, accessTokenTTL, refreshTokenTTL)
	if err != nil {
		log.Fatalln(err)
	}

	services := service.NewServices(service.ServicesDeps{
		Repos:           repos,
		Hasher:          hasher,
		TokenManager:    tokenManager,
		AccessTokenTTL:  accessTokenTTL,
		RefreshTokenTTL: refreshTokenTTL,
	})

	go services.Users.DeleteExpiredSessions()

	router := gorouter.NewRouter()

	handler := handler.NewHandler(services.Users, services.Posts, services.Comments, tokenManager)
	handler.Init(router)

	server := server.NewServer(config, router)
	log.Fatalln(server.Run())
}
