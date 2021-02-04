package handlers

import "net/http"

func APIHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"api":"works"}`))
}
