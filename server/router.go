package server

import (
	"labyrinth/server/handlers"
	"labyrinth/server/handlers/auth"
	"labyrinth/server/handlers/company"
	"labyrinth/server/handlers/user"

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
└── user/ # POST
    ├── {user_id}/ # GET, POST, DELETE
    │   │  └── profile # GET, POST, DELETE
    │   │
    │   └── company/ # GET, POST
    │       └──  {company_id}/ # GET
	│				   ├── profile # GET, POST,  DELETE
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

	// проверка сервера на готовность
	r.HandleFunc("labyrinth/ping", handlers.Ping).Methods("GET")

	// авторизация
	r.HandleFunc("labyrinth/auth/register", auth.RegisterUserHandler).Methods("POST")
	r.HandleFunc("labyrinth/auth/login", auth.LoginUserHandler).Methods("POST")
	// r.HandleFunc("labyrinth/auth/reset", handler).Methods("POST")

	// работа с пользователем
	r.HandleFunc("labyrinth/user/{user_id}/profile", user.GetUserProfileHandler).Methods("GET")
	r.HandleFunc("labyrinth/user/{user_id}/profile", user.UpdateUserProfileHandler).Methods("POST")
	// r.HandleFunc("labyrinth/user/{user_id}/profile", user.DeleteUserProfileHandler).Methods("DLETE")

	// работа с компанией
	r.HandleFunc("labyrinth/user/{user_id}/company", company.NewCompanyHandler).Methods("POST")
	r.HandleFunc("labyrinth/user/{user_id}/company", company.GetAllCompaniesHandler).Methods("GET")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}", company.GetCompanyHandler).Methods("GET")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/profile", company.GetCompanyProfileHandler).Methods("GET")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/profile", company.UpdateCompanyProfileHandler).Methods("POST")
	// r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/profile", company.DeletCompanyProfileHandler).Methods("DELETE")

	// r.HandleFunc("labyrinth/user/{user_id}/compnay/{company_id}/invite", handler).Methods("GET", "POST", "DELETE")

	// r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}", handler).Methods("GET") где get - вход в комапию
	// r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/profile", handler).Methods("GET", "POST")

	// r.HandleFunc("labyrinth/user/{user_id}/compnay/{company_id}/employee", handler).Methods("GET", "POST")
	// r.HandleFunc("labyrinth/user/{user_id}/compnay/{company_id}/employee/{employee_id}", handler).Methods("GET", "PUT", "DELETE")

	// r.HandleFunc("labyrinth/user/{user_id}/company/{compnay_id}/department", handler).Methods("GET", "POST")
	// r.HandleFunc("labyrinth/user/{user_id}/company/{compnay_id}/department/{department_id}", handler).Methods("GET", "PUT", "DELETE")

	// r.HandleFunc("labyrinth/user/{user_id}/company/{compnay_id}/department/{department_id}/notebook").Method("GET", "POST")
	// r.HandleFunc("labyrinth/user/{user_id}/company/{compnay_id}/department/{department_id}/notebook/{notebook_id}").Method("GET", "POST")

	return r
}
