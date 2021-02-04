package api

import (
	"net/http"

	"github.com/alseiitov/real-time-forum/internal/api/handlers"
	"github.com/alseiitov/real-time-forum/internal/api/middlewares"
)

func Init() {
	http.HandleFunc("/api/user", middlewares.EnableCORS(handlers.UserHandler))
}
