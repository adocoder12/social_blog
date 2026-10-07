package handler

import "net/http"

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the site health\n"))

}
