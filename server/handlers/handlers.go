package handlers

import (
	"labyrinth/server/handlers/auth"
	"labyrinth/server/handlers/company"
	"labyrinth/server/handlers/department"
	"labyrinth/server/handlers/depemployee"
	"labyrinth/server/handlers/depposition"
	"labyrinth/server/handlers/employee"
	"labyrinth/server/handlers/journal"
	"labyrinth/server/handlers/permission"
	"labyrinth/server/handlers/position"
	"labyrinth/server/handlers/user"
	"net/http"
)

type authInterface interface {
	LoginUserHandler(w http.ResponseWriter, r *http.Request)
	RegisterUserHandler(w http.ResponseWriter, r *http.Request)
}

type userInterface interface {
	GetUserProfileHandler(w http.ResponseWriter, r *http.Request)
	UpdateUserProfileHandler(w http.ResponseWriter, r *http.Request)
}

type companyInterface interface {
	GetAllCompaniesHandler(w http.ResponseWriter, r *http.Request)
	GetCompanyHandler(w http.ResponseWriter, r *http.Request)
	NewCompanyHandler(w http.ResponseWriter, r *http.Request)
	UpdateCompanyProfileHandler(w http.ResponseWriter, r *http.Request)
	GetCompanyProfileHandler(w http.ResponseWriter, r *http.Request)
}

type employeeInterface interface {
	GetAllEmployeeHandler(w http.ResponseWriter, r *http.Request)
	NewEmployeeHandler(w http.ResponseWriter, r *http.Request)
	UpdateEmployeeHandler(w http.ResponseWriter, r *http.Request)
}

type positionInterface interface {
	GetAllPositionHandler(w http.ResponseWriter, r *http.Request)
	NewPositionHandler(w http.ResponseWriter, r *http.Request)
	UpdatePositionHandler(w http.ResponseWriter, r *http.Request)
}

type departmentInterface interface {
	GetDepartmentProfileHandler(w http.ResponseWriter, r *http.Request)
	GetDepartmentHandler(w http.ResponseWriter, r *http.Request)
	NewDepartmentHandler(w http.ResponseWriter, r *http.Request)
	UpdateDepartmentProfileHandler(w http.ResponseWriter, r *http.Request)
	UpdateDepartmentHandler(w http.ResponseWriter, r *http.Request)
}

type depemployeeInterface interface {
	GetAllDepEmployeeHandler(w http.ResponseWriter, r *http.Request)
	NewDepEmployeeHandler(w http.ResponseWriter, r *http.Request)
	UpdateDepEmployeeHandler(w http.ResponseWriter, r *http.Request)
}

type depemployeePosInterface interface {
	GetAllDepPositionHandler(w http.ResponseWriter, r *http.Request)
	NewDepPositionHandler(w http.ResponseWriter, r *http.Request)
	UpdateDepPositionHandler(w http.ResponseWriter, r *http.Request)
}

type notebookInterface interface {
	NewNotebookHandler(w http.ResponseWriter, r *http.Request)
	GetNotebookHandler(w http.ResponseWriter, r *http.Request)
	UpdateNotebookHandler(w http.ResponseWriter, r *http.Request)
}

type permissionInterface interface {
	GetPermissionHandler(w http.ResponseWriter, r *http.Request)
	UpdatePermissionHandler(w http.ResponseWriter, r *http.Request)
}

type Handlers struct {
	Auth                       authInterface
	UserProfile                userInterface
	Company                    companyInterface
	Employee                   employeeInterface
	Position                   positionInterface
	Department                 departmentInterface
	DepartmentEmployee         depemployeeInterface
	DepartmentEmployeePosition depemployeePosInterface
	Notebook                   notebookInterface
	Permission                 permissionInterface
}

func NewHandlers() Handlers {
	return Handlers{
		Auth:                       auth.NewAuthHandlers(),
		UserProfile:                user.NewUserHandlers(),
		Company:                    company.NewCompanyHandlers(),
		Employee:                   employee.NewEmployeeHandlers(),
		Position:                   position.NewPositionHandlers(),
		Department:                 department.NewDepartmentHandlers(),
		DepartmentEmployee:         depemployee.NewDepEmployeeHandlers(),
		DepartmentEmployeePosition: depposition.NewDepPositionHandlers(),
		Notebook:                   journal.NewJournalHandler(),
		Permission:                 permission.NewPermissionHandlers(),
	}
}
