package main

import (
	"flag"

	"github.com/alseiitov/real-time-forum/internal/server"
)

func main() {
	configPath := flag.String("config-path", "./configs/config.json", "Path to the config file")
	flag.Parse()

	server.Run(configPath)
}
