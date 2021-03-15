package main

import (
	"flag"
	"log"
	"net/http"
	"text/template"

	"github.com/alseiitov/real-time-forum/internal/config"
)

func main() {
	var configPath = flag.String("config-path", "./configs/config.json", "Path to the config file")
	flag.Parse()

	conf, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatalln(err)
	}

	temp, err := template.ParseFiles("./web/public/index.html")
	if err != nil {
		log.Fatalln(err)
	}

	fileServer := http.FileServer(http.Dir("./web/src"))

	http.Handle("/src/", http.StripPrefix("/src/", fileServer))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := temp.Execute(w, conf.BackendAdress())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	})

	port := conf.FrontendPort()
	log.Printf("Frontend server is starting at %v", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalln(err)
	}
}
