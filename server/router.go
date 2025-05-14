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
	r.HandleFunc("labyrinth/user/{user_id}/profile", manager.UserProfile.GetUserProfileHandler).Methods("GET")
	r.HandleFunc("labyrinth/user/{user_id}/profile", manager.UserProfile.UpdateUserProfileHandler).Methods("POST")
	// r.HandleFunc("labyrinth/user/{user_id}/profile", user.DeleteUserProfileHandler).Methods("DLETE")

	// работа с компанией
	r.HandleFunc("labyrinth/user/{user_id}/company", manager.Company.NewCompanyHandler).Methods("POST")
	r.HandleFunc("labyrinth/user/{user_id}/company", manager.Company.GetAllCompaniesHandler).Methods("GET")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}", manager.Company.GetCompanyHandler).Methods("GET")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/profile", manager.Company.GetCompanyProfileHandler).Methods("GET")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/profile", manager.Company.UpdateCompanyProfileHandler).Methods("POST")
	// r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/profile", company.DeletCompanyProfileHandler).Methods("DELETE")

	// работа с позициями
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/position", manager.Position.GetAllPositionHandler).Methods("GET")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/position", manager.Position.NewPositionHandler).Methods("POST")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/position/{position_id}", manager.Position.UpdatePositionHandler).Methods("POST")

	// работа с инвайтами
	// r.HandleFunc("labyrinth/user/{user_id}/compnay/{company_id}/invite", handler).Methods("GET", "POST", "DELETE")

	// работа с работниками
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/employee", manager.Employee.GetAllEmployeeHandler).Methods("GET")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/employee", manager.Employee.NewEmployeeHandler).Methods("POST")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/employee/{employee_id}", manager.Employee.UpdateEmployeeHandler).Methods("POST")
	// r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/employee/{employee_id}", employee.DeleteEmployeeHandler).Methods("DELETE")

	// работа с департаментами
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department", manager.Department.NewDepartmentHandler).Methods("POST")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department/{department_id}", manager.Department.GetDepartmentHandler).Methods("GET")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department/{department_id}", manager.Department.UpdateDepartmentHandler).Methods("POST")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department/{department_id}/profile", manager.Department.GetDepartmentProfileHandler).Methods("GET")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department/{department_id}/profile", manager.Department.UpdateDepartmentProfileHandler).Methods("POST")
	// r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department/{department_id}/profile", department.DeleteDepartmentProfileHandler).Methods("DELETE")

	// работа с работниками департаментов
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department/{department_id}/depemployee", manager.DepartmentEmployee.GetAllDepEmployeeHandler).Methods("POST")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department/{department_id}/depemployee", manager.DepartmentEmployee.NewDepEmployeeHandler).Methods("POST")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department/{department_id}/depemployee/{depemployee_id}", manager.DepartmentEmployee.UpdateDepEmployeeHandler).Methods("POST")
	// r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department/{department_id}/depemployee/{depemployee_id}", depemployee.DeleteDepEmployeeHandler).Methods("DELETE")

	// работа с позициями работников департаментов
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department/{department_id}/depposition", manager.DepartmentEmployeePosition.GetAllDepPositionHandler).Methods("GET")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department/{department_id}/depposition", manager.DepartmentEmployeePosition.NewDepPositionHandler).Methods("POST")
	r.HandleFunc("labyrinth/user/{user_id}/company/{company_id}/department/{department_id}/depposition/{depposition_id}", manager.DepartmentEmployeePosition.UpdateDepPositionHandler).Methods("POST")

	// работа с лабораторными  журналами
	// r.HandleFunc("labyrinth/user/{user_id}/company/{compnay_id}/department/{department_id}/notebook").Method("GET", "POST")
	// r.HandleFunc("labyrinth/user/{user_id}/company/{compnay_id}/department/{department_id}/notebook/{notebook_id}").Method("GET", "POST")

	return r
}
