package server

import (
	"labyrinth/server/handlers"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/ping", handlers.Ping)

	return r
}
