package server

import (
	"labyrinth/server/handlers"
	"labyrinth/server/middleware"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("app/ping", handlers.Ping).Methods("GET")

	r.HandleFunc("app/register", handlers.RegisterUserHandler).Methods("POST")
	r.HandleFunc("app/login", handlers.LoginUserHandler).Methods("POST")

	r.HandleFunc("app/companies", middleware.AuthMiddleware(handlers.GetCompaniesHandler)).Methods("GET")
	r.HandleFunc("app/company/{id}", middleware.AuthMiddleware(handlers.LoginCompanyHandler)).Methods("GET")
	return r
}
