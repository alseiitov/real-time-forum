package main

import (
	"flag"

	"github.com/alseiitov/real-time-forum/internal/app"
)

func main() {
	configPath := flag.String("config-path", "./configs/config.json", "Path to the config file")
	flag.Parse()

	app.Run(configPath)
}
