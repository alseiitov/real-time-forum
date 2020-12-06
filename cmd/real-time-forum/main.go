package main

import (
	"log"
	"net/http"
)

func main() {
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln(err)
	}
}
