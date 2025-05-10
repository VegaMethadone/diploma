package server

import (
	"labyrinth/server/handlers"

	"github.com/gorilla/mux"
)

/*

labyrinth/
├── auth/
│   ├── register # POST
│   ├── login # POST
│   └── reset # POST
│
├── ping # GET
│
│
└── user/
    ├── {user_id}/ # GET, POST, DELETE
    │   │  └── profile # GET, POST
    │   │
    │   └── company/ # GET, POST
    │       └──  {company_id}/ # GET, POST, DELETE
	│				   ├── profile # GET, POST
	│				   ├── invite  # GET, POST
	│				   ├──	employee/  # GET, POST
	│				   │ 		└── {employee_id}   # GET, POST, DELETE
	│				   │
    │                  ├── department/ # GET, POST
    │                  │   └── {department_id} # GET, POST, DELETE
	│								 ├── profile # GET, POST
    │                  │             │
	│                  │             └── depemployee/ # GET, POST
	│		           │                      └──{depemployee_id} # GET, POST, PUT, DELETE
	│			       │
	│			       │
    │                  └── notebook/ # GET, POST
    │                      └── {notebook_id} # GET, POST, DELETE
    │
    └── тут будет онбординг ?

*/

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("labyrinth/ping", handlers.Ping).Methods("GET")

	// r.HandleFunc("labyrinth/auth/register", handler).Methods("POST")
	// r.HandleFunc("labyrinth/auth/login", handler).Methods("POST")
	// r.HandleFunc("labyrinth/auth/reset", handler).Methods("POST")

	// r.HandleFunc("labyrinth/user/{user_id}", handler).Methods("GET")
	// r.HandleFunc("labyrinth/user/{user_id}/profile", handler).Methods("GET", "POST")
	// r.HandleFunc("labyrinth/user/{user_id}/compnay", handler).Methods("GET", "POST")

	// r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}", handler).Methods("GET") где get - вход в комапию
	// r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/profile", handler).Methods("GET", "POST")
	// r.HandleFunc("labyrinth/user/{user_id}/compnay/{company_id}/invite", handler).Methods("GET", "POST")

	// r.HandleFunc("labyrinth/user/{user_id}/compnay/{company_id}/employee", handler).Methods("GET", "POST")
	// r.HandleFunc("labyrinth/user/{user_id}/compnay/{company_id}/employee/{employee_id}", handler).Methods("GET", "PUT", "DELETE")

	// r.HandleFunc("labyrinth/user/{user_id}/company/{compnay_id}/department", handler).Methods("GET", "POST")
	// r.HandleFunc("labyrinth/user/{user_id}/company/{compnay_id}/department/{department_id}", handler).Methods("GET", "PUT", "DELETE")

	// r.HandleFunc("labyrinth/user/{user_id}/company/{compnay_id}/department/{department_id}/notebook").Method("GET", "POST")
	// r.HandleFunc("labyrinth/user/{user_id}/company/{compnay_id}/department/{department_id}/notebook/{notebook_id}").Method("GET", "POST")

	return r
}
