package logic

import (
	authlogic "labyrinth/logic/authLogic"
	companylogic "labyrinth/logic/companyLogic"
	departmentlogic "labyrinth/logic/departmentLogic"
	depemployeelogic "labyrinth/logic/depemployeeLogic"
	depemployeeposlogic "labyrinth/logic/depemployeeposLogic"
	employeelogic "labyrinth/logic/employeeLogic"
	positionlogic "labyrinth/logic/positionLogic"
	userlogic "labyrinth/logic/userLogic"
	"labyrinth/models/company"
	"labyrinth/models/department"
	"labyrinth/models/depemployee"
	"labyrinth/models/depposition"
	"labyrinth/models/employee"
	"labyrinth/models/position"
	"labyrinth/models/user"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type authLogic interface {
	Login(mail, password string) (*user.User, error)
	Register(mail, hashPassword, phone string) error
}

type userLogic interface {
	UpdateUserProfile(userProfile *user.User) error
	GetUserProfile(userId uuid.UUID) (*user.User, error)
}

type companyLogic interface {
	DeleteCompany(companyId uuid.UUID) error
	GetCompany(userId, companyId uuid.UUID) (*company.Company, error)
	GetUserCompanies(userId uuid.UUID) (*[]company.Company, error)
	NewCompany(userId uuid.UUID, name, description string) (uuid.UUID, error)
	UpdateCompany(comp *company.Company, companyId, employeeId uuid.UUID) error
}

type employeeLogic interface {
	DeleteEmployee(userId, companyId, employeeId uuid.UUID) error
	GetAllEmployee(companyId uuid.UUID) (*[]employee.Employee, error)
	GetEmployee(userId, companyId uuid.UUID) (*employee.Employee, error)
	NewEmployee(employeeId, userId, companyId, positionId uuid.UUID) error
	UpdateEmployee(userId, companyId uuid.UUID, updatedEmployee *employee.Employee) error
}

type positionLogic interface {
	GetAllPositions(userId, companyId uuid.UUID) (*[]position.Position, error)
	NewPosition(userId, companyId uuid.UUID, lvl int, name string) (uuid.UUID, error)
	UpdatePosition(userId, companyId uuid.UUID, updatePosition *position.Position) error
}

type departmentLogic interface {
	DeleteDepartment(userId, companyId, departmentId uuid.UUID) error
	GetDepartment(userId, companyId, departmentId uuid.UUID) (*department.Department, error)
	NewDepartment(userId, companyId, parentId uuid.UUID, name, description string) (uuid.UUID, uuid.UUID, uuid.UUID, error)
	UpdateDepartment(userId, companyId uuid.UUID, updateDepartment *department.Department) error
}

type departmentEmployeeLogic interface {
	DeleteDepartmentEmployee(employeeId, departmentId, depemployeeId uuid.UUID) error
	GetAllDepEmployees(departmentId uuid.UUID) (*[]depemployee.DepartmentEmployee, error)
	GetDepartmentEmployee(employeeId, departmentId uuid.UUID) (*depemployee.DepartmentEmployee, error)
	NewDepemployee(employeeId, departmentId, positionId uuid.UUID) error
	UpdateDepEmployee(employeeId, departmentId uuid.UUID, updatedDepEmployee *depemployee.DepartmentEmployee) error
}

type departmentEmployeePosLogic interface {
	DeleteDepEmployeePos(currentlvl int, employeeId, departmentId, positionId uuid.UUID) error
	GetAllDepEmployeePos(departmentId uuid.UUID) (*[]depposition.DepPosition, error)
	NewDepemployeePos(departmentId uuid.UUID, lvl int, name string) (uuid.UUID, error)
	UpdateDepEmployeePos(currentlvl int, employeeId, departmentId uuid.UUID, position *depposition.DepPosition) error
}

type jwtLogic interface {
	NewToken(settings jwt.MapClaims) string
	VerifyToken(tokenString string) (jwt.MapClaims, error)
}

type BusinessLogic struct {
	Jwt                        jwtLogic
	Auth                       authLogic
	User                       userLogic
	Company                    companyLogic
	Position                   positionLogic
	Employee                   employeeLogic
	Department                 departmentLogic
	DepartmentEmployee         departmentEmployeeLogic
	DepartmentEmployeePosition departmentEmployeePosLogic
}

func NewBusinessLogic() *BusinessLogic {
	return &BusinessLogic{
		Jwt:                        NewMyJwt(),
		Auth:                       authlogic.NewAuth(),
		User:                       userlogic.NewUserlogic(),
		Company:                    companylogic.NewCompanyLogic(),
		Position:                   positionlogic.NewPositionLogic(),
		Employee:                   employeelogic.NewEmployeeLogic(),
		Department:                 departmentlogic.NewDepartmentLogic(),
		DepartmentEmployee:         depemployeelogic.NewDepemployeeLogic(),
		DepartmentEmployeePosition: depemployeeposlogic.NewDepemploeePosLogic(),
	}
}
