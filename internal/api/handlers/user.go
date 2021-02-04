package handlers

import "net/http"

func UserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"api":"works"}`))
}
