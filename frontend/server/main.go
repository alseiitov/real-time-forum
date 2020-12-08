package main

import (
	"log"
	"flag"
	"net/http"

	"github.com/alseiitov/real-time-forum/frontend/server/configs"
)

func main() {
	var configPath = flag.String("config-path", "./configs/config.json", "Path to the config file")
	flag.Parse()

	conf, err := configs.Read(*configPath)
	if err != nil {
		log.Fatalln(err)
	}

	fileServer := http.FileServer(http.Dir("./src"))
	http.Handle("/", fileServer)

	port := conf.Server.Port
	log.Printf("Frontend server is starting at port %v", port)
	if err := http.ListenAndServe(":" + port, nil); err != nil {
		log.Fatalln(err)
	}
}
