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
	hasher := hash.NewBcryptHasher()

	secretKey := os.Getenv("SECRET_KEY")
	accessTokenTTL, refreshTokenTTL, err := config.GetTokenTTLs()
	if err != nil {
		log.Fatalln(err)
	}

	tokenManager, err := auth.NewManager(secretKey, accessTokenTTL, refreshTokenTTL)
	if err != nil {
		log.Fatalln(err)
	}

	services := service.NewServices(service.ServicesDeps{
		Repos:        repos,
		Hasher:       hasher,
		TokenManager: tokenManager,
	})

	router := gorouter.NewRouter()

	handler := handler.NewHandler(services.Users)
	handler.Init(router)

	server := server.NewServer(config, router)
	server.Run()
}
