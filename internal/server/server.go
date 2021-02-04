package server

import (
	"log"
	"net/http"

	"github.com/alseiitov/real-time-forum/internal/api"
	"github.com/alseiitov/real-time-forum/internal/configs"
	"github.com/alseiitov/real-time-forum/internal/storage"
)

func Run(configPath *string) {
	config, err := configs.NewConfig(*configPath)
	if err != nil {
		log.Fatalln(err)
	}

	db, err := storage.ConnectDB(config)
	defer db.Close()
	if err != nil {
		log.Fatalln(err)
	}

	port := config.GetBackendPort()

	api.Init()

	log.Printf("Backend server is starting at port %v", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalln(err)
	}

}
