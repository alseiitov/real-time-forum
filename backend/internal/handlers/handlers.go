package handlers

import (
	"net/http"

	"github.com/alseiitov/real-time-forum/backend/internal/middlewares"
)

func Init() {
	http.HandleFunc("/api", middlewares.EnableCORS(APIHandler))
}
