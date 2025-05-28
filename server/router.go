package server

import (
	"labyrinth/server/handlers"
	"labyrinth/server/middleware"

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
	manager := handlers.NewHandlers()

	// проверка сервера на готовность
	r.HandleFunc("labyrinth/ping", handlers.Ping).Methods("GET")

	// авторизация
	r.HandleFunc("labyrinth/auth/register", manager.Auth.RegisterUserHandler).Methods("POST")
	r.HandleFunc("labyrinth/auth/login", manager.Auth.LoginUserHandler).Methods("POST")
	// r.HandleFunc("labyrinth/auth/reset", handler).Methods("POST")

	// работа с пользователем
	r.HandleFunc("labyrinth/user/{user_id}/profile", middleware.AuthMiddleware(manager.UserProfile.GetUserProfileHandler)).Methods("GET")
	r.HandleFunc("labyrinth/user/{user_id}/profile", middleware.AuthMiddleware(manager.UserProfile.UpdateUserProfileHandler)).Methods("POST")
	// r.HandleFunc("labyrinth/user/{user_id}/profile", user.DeleteUserProfileHandler).Methods("DLETE")

	// работа с компанией
	r.HandleFunc("labyrinth/user/{user_id}/company", middleware.AuthMiddleware(manager.Company.NewCompanyHandler)).Methods("POST")
	r.HandleFunc("labyrinth/user/{user_id}/company", middleware.AuthMiddleware(manager.Company.GetAllCompaniesHandler)).Methods("GET")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}", middleware.AuthMiddleware(manager.Company.GetCompanyHandler)).Methods("GET")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/profile", middleware.AuthMiddleware(manager.Company.GetCompanyProfileHandler)).Methods("GET")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/profile", middleware.AuthMiddleware(manager.Company.UpdateCompanyProfileHandler)).Methods("POST")
	// r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/profile", company.DeletCompanyProfileHandler).Methods("DELETE")

	// работа с позициями
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/position", middleware.AuthMiddleware(manager.Position.GetAllPositionHandler)).Methods("GET")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/position", middleware.AuthMiddleware(manager.Position.NewPositionHandler)).Methods("POST")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/position/{position_id}", middleware.AuthMiddleware(manager.Position.UpdatePositionHandler)).Methods("POST")

	// работа с инвайтами
	// r.HandleFunc("labyrinth/user/{user_id}/compnay/{company_id}/invite", handler).Methods("GET")
	// r.HandleFunc("labyrinth/user/{user_id}/compnay/{company_id}/invite", handler).Methods("POST")
	// r.HandleFunc("labyrinth/user/{user_id}/compnay/{company_id}/invite/{invite_id}", handler).Methods("DELETE")

	// работа с работниками
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/employee", middleware.AuthMiddleware(manager.Employee.GetAllEmployeeHandler)).Methods("GET")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/employee", middleware.AuthMiddleware(manager.Employee.NewEmployeeHandler)).Methods("POST")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/employee/{employee_id}", middleware.AuthMiddleware(manager.Employee.UpdateEmployeeHandler)).Methods("POST")
	// r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/employee/{employee_id}", employee.DeleteEmployeeHandler).Methods("DELETE")

	// работа с департаментами
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department", middleware.AuthMiddleware(manager.Department.NewDepartmentHandler)).Methods("POST")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department/{department_id}", middleware.AuthMiddleware(manager.Department.GetDepartmentHandler)).Methods("GET")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department/{department_id}", middleware.AuthMiddleware(manager.Department.UpdateDepartmentHandler)).Methods("POST")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department/{department_id}/profile", middleware.AuthMiddleware(manager.Department.GetDepartmentProfileHandler)).Methods("GET")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department/{department_id}/profile", middleware.AuthMiddleware(manager.Department.UpdateDepartmentProfileHandler)).Methods("POST")
	// r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department/{department_id}/profile", department.DeleteDepartmentProfileHandler).Methods("DELETE")

	// работа с работниками департаментов
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department/{department_id}/depemployee", middleware.AuthMiddleware(manager.DepartmentEmployee.GetAllDepEmployeeHandler)).Methods("POST")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department/{department_id}/depemployee", middleware.AuthMiddleware(manager.DepartmentEmployee.NewDepEmployeeHandler)).Methods("POST")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department/{department_id}/depemployee/{depemployee_id}", middleware.AuthMiddleware(manager.DepartmentEmployee.UpdateDepEmployeeHandler)).Methods("POST")
	// r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department/{department_id}/depemployee/{depemployee_id}", depemployee.DeleteDepEmployeeHandler).Methods("DELETE")

	// работа с позициями работников департаментов
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department/{department_id}/depposition", middleware.AuthMiddleware(manager.DepartmentEmployeePosition.GetAllDepPositionHandler)).Methods("GET")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department/{department_id}/depposition", middleware.AuthMiddleware(manager.DepartmentEmployeePosition.NewDepPositionHandler)).Methods("POST")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department/{department_id}/depposition/{depposition_id}", middleware.AuthMiddleware(manager.DepartmentEmployeePosition.UpdateDepPositionHandler)).Methods("POST")

	// работа с лабораторными  журналами
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department/{department_id}/notebook", middleware.AuthMiddleware(manager.Notebook.NewNotebookHandler)).Methods("POST")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department/{department_id}/notebook/{notebook_id}", middleware.AuthMiddleware(manager.Notebook.GetNotebookHandler)).Methods("GET")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department/{department_id}/notebook/{notebook_id}", middleware.AuthMiddleware(manager.Notebook.UpdateNotebookHandler)).Methods("POST")

	// работа с разрешениями журнала
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department/{department_id}/notebook/{notebook_id}/permission", middleware.AuthMiddleware(manager.Permission.GetPermissionHandler)).Methods("GET")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department/{department_id}/notebook/{notebook_id}/permission", middleware.AuthMiddleware(manager.Permission.UpdatePermissionHandler)).Methods("POST")
	return r
}
