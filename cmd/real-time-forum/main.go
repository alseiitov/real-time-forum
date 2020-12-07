package main

import (
	"log"
	"flag"

	"github.com/alseiitov/real-time-forum/configs"
	"github.com/alseiitov/real-time-forum/internal/server"
)

func main() {
	var configPath = flag.String("config-path", "./configs/config.json", "Path to the config file")
	flag.Parse()
	
	config, err := configs.Read(*configPath)
	if err != nil {
		log.Fatalln(err)
	}

	server.Run(config)
}
