package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"labyrinth/config"
	"labyrinth/models/company"
	"labyrinth/models/employee"
	"labyrinth/models/position"
	"labyrinth/models/user"

	"github.com/google/uuid"
)

type userDB interface {
	// CreateUser создает нового пользователя (возвращает созданного пользователя с заполненными полями)
	CreateUser(
		ctx context.Context,
		// db *sql.DB,
		sharedTx *sql.Tx,
		u *user.User,
	) (*user.User, error)

	// GetUserByCredentials ищет пользователя по логину и проверяет пароль
	GetUserByCredentials(
		ctx context.Context,
		// db *sql.DB,
		sharedTx *sql.Tx,
		login,
		password string,
	) (*user.User, error)

	// GetUserByID получает пользователя по ID (без проверки пароля)
	GetUserByID(
		ctx context.Context,
		// db *sql.DB,
		sharedTx *sql.Tx,
		id uuid.UUID,
	) (*user.User, error)

	// UpdateUser обновляет данные пользователя
	UpdateUser(
		ctx context.Context,
		// db *sql.DB,
		sharedTx *sql.Tx,
		u *user.User,
	) error

	// DeleteUser мягкое удаление (is_active = false)
	DeleteUser(
		ctx context.Context,
		// db *sql.DB,
		sharedTx *sql.Tx,
		id uuid.UUID,
	) error
}

type companyDB interface {
	// CreateCompany создает новую компанию
	CreateCompany(
		ctx context.Context,
		sharedTx *sql.Tx,
		company *company.Company,
	) error

	// GetCompanyByID получает компанию по ID
	GetCompanyByID(
		ctx context.Context,
		sharedTx *sql.Tx,
		id uuid.UUID,
	) (*company.Company, error)

	// AddUserToCompany добавление юзера в компанию
	AddUserToCompany(
		ctx context.Context,
		sharedTx *sql.Tx,
		userID uuid.UUID,
		companyID uuid.UUID,
	) error

	// GetCompaniesByUser получает компании пользователя
	GetCompaniesByUser(
		ctx context.Context,
		sharedTx *sql.Tx,
		userID uuid.UUID,
	) ([]*company.Company, error)

	// UpdateCompany обновляет данные компании
	UpdateCompany(
		ctx context.Context,
		sharedTx *sql.Tx,
		company *company.Company,
	) error

	// DeleteCompany помечает компанию как удаленную
	DeleteCompany(
		ctx context.Context,
		sharedTx *sql.Tx,
		id uuid.UUID,
	) error

	// DeactivateCompanyUsers деактивирует компанию у юзера
	DeactivateCompanyUsers(
		ctx context.Context,
		sharedTx *sql.Tx,
		companyID uuid.UUID,
	) error
}

type employeeDB interface {
	// CreateEmployee создает нового сотрудника в базе данных.
	CreateEmployee(
		ctx context.Context,
		sharedTx *sql.Tx,
		empl *employee.Employee,
	) error

	// UpdateEmployee обновляет данные сотрудника.
	UpdateEmployee(
		ctx context.Context,
		sharedTx *sql.Tx,
		empl *employee.Employee,
	) error

	// GetEmployee возвращает сотрудника по его ID.
	GetEmployeeByUserId(
		ctx context.Context,
		sharedTx *sql.Tx,
		employeeId uuid.UUID,
	) (*employee.Employee, error)

	// GetEmployeesByCompanyId возвращает список сотрудников компании.
	GetEmployeesByCompanyId(
		ctx context.Context,
		sharedTx *sql.Tx,
		companyId uuid.UUID,
	) ([]*employee.Employee, error)

	// DeleteEmployee is_active = false
	DeleteEmployee(
		ctx context.Context,
		sharedTx *sql.Tx,
		employeeId uuid.UUID,
	) error

	// ExistsEmployee проверяет, существует ли сотрудник с таким ID.
	ExistsEmployee(
		ctx context.Context,
		sharedTx *sql.Tx,
		employeeId uuid.UUID,
	) (bool, error)

	// CountEmployees возвращает количество сотрудников в компании.
	CountEmployees(
		ctx context.Context,
		sharedTx *sql.Tx,
		companyId uuid.UUID,
	) (int, error)
}

type positionDB interface {
	// CreatePosition создает новую должность в компании
	CreatePosition(
		ctx context.Context,
		sharedTx *sql.Tx,
		position *position.Position,
	) error

	// GetPositionById возвращает должность по UUID
	GetPositionById(
		ctx context.Context,
		sharedTx *sql.Tx,
		positionId uuid.UUID,
	) (*position.Position, error)

	// GetPositionsByCompanyId возвращает все активные должности компании
	GetPositionsByCompanyId(
		ctx context.Context,
		sharedTx *sql.Tx,
		companyId uuid.UUID,
	) ([]*position.Position, error)

	// UpdatePosition обновляет данные должности
	UpdatePosition(
		ctx context.Context,
		sharedTx *sql.Tx,
		position *position.Position,
	) error

	// DeletePosition мягко удаляет должность (is_active = false)
	DeletePosition(
		ctx context.Context,
		sharedTx *sql.Tx,
		positionId uuid.UUID,
	) error
}
type departmentDB interface {
	// GetEmployeesByDepartment возвращает сотрудников по отделу.
	GetEmployeesByDepartment(
		ctx context.Context,
		sharedTx *sql.Tx,
		departmentId uuid.UUID,
	) ([]*employee.Employee, error)
}
type departmentEmployeeDB interface{}

type departmentEmployeePositionDB interface{}

type PostgresDB struct {
	User                       userDB
	Company                    companyDB
	Employee                   employeeDB
	Position                   positionDB
	Department                 departmentDB
	DepartmentEmployee         departmentEmployeeDB
	DepartmentEmployeePosition departmentEmployeePositionDB
}

func GetConnection() string {
	conf := config.Conf.PostgreSQL
	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=%s host=%s",
		conf.Username,
		conf.Password,
		conf.DatabaseName,
		conf.SSLMode,
		conf.Host,
	)
	return connStr
}
