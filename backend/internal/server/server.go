package server

import (
	"log"
	"net/http"

	"github.com/alseiitov/real-time-forum/backend/internal/configs"
	"github.com/alseiitov/real-time-forum/backend/internal/handlers"
)

func Run(config *configs.Conf) {
	handlers.HandleFunctions()
	port := config.GetPort()

	log.Printf("Backend server is starting at port %v", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalln(err)
	}
}
