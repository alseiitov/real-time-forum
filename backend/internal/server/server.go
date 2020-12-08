package server

import (
	"net/http"
	"log"
	"github.com/alseiitov/real-time-forum/backend/internal/configs"
)

func Run(config *configs.Conf)  {
	http.HandleFunc("/api", apiHandler)

	port := config.GetPort()
	log.Printf("Backend server is starting at port %v", port)
	if err := http.ListenAndServe(":" + port, nil); err != nil {
		log.Fatalln(err)
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.Write([]byte(`{"api":"works"}`))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
 }

