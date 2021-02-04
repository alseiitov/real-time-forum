package handlers

import (
	"net/http"

	"github.com/alseiitov/real-time-forum/internal/middlewares"
)

func Init() {
	http.HandleFunc("/api", middlewares.EnableCORS(APIHandler))
}
