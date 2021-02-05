package app

import (
	"log"

	"github.com/alseiitov/real-time-forum/internal/config"
	"github.com/alseiitov/real-time-forum/internal/repository"
	"github.com/alseiitov/real-time-forum/internal/service"
	"github.com/alseiitov/real-time-forum/pkg/database/sqlite"
)

func Run(configPath string) {
	config, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	db, err := sqlite.ConnectDB(config)
	if err != nil {
		log.Fatalln(err)
	}

	repos := repository.NewRepositories(db)

	services := service.NewServices(service.ServicesDeps{
		Repos: repos,
	})

}
