package server

import (
	"labyrinth/server/handlers"
	"labyrinth/server/middleware"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/ping", handlers.Ping).Methods("GET")

	r.HandleFunc("/register", handlers.RegisterUserHandler).Methods("POST")
	r.HandleFunc("/login", handlers.LoginUserHandler).Methods("POST")

	r.HandleFunc("/companies", middleware.AuthMiddleware(handlers.GetCompaniesHandler)).Methods("GET") // посмотреть, мог накосяить с ограничем методов
	return r
}
