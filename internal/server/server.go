package server

import (
	"net/http"
	"log"
	"github.com/alseiitov/real-time-forum/configs"
)

func Run(config *configs.Conf)  {
	http.HandleFunc("/", indexHandler)

	port := config.GetPort()
	log.Printf("Server is starting at port %v", port)
	if err := http.ListenAndServe(":" + port, nil); err != nil {
		log.Fatalln(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}