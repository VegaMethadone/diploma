package server

import (
	"labyrinth/server/handlers"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/ping", handlers.Ping).Methods("GET")
	/*
		r.HandleFunc("/register")
		r.HandleFunc("/login")
		r.HandleFunc("/profile/{id}").Methods("GET", "PUT", "DELETE", "POST")
	*/
	return r
}
