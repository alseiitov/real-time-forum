package main

import (
	"flag"
	"log"

	"github.com/alseiitov/real-time-forum/backend/internal/configs"
	"github.com/alseiitov/real-time-forum/backend/internal/server"
	"github.com/alseiitov/real-time-forum/backend/internal/storage"
)

func main() {
	configPath := flag.String("config-path", "./configs/config.json", "Path to the config file")
	flag.Parse()

	config, err := configs.NewConfig(*configPath)
	if err != nil {
		log.Fatalln(err)
	}

	db, err := storage.ConnectDB(config)
	defer db.Close()
	if err != nil {
		log.Fatalln(err)
	}

	server.Run(config, db)
}
