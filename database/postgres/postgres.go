package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"labyrinth/config"
	"labyrinth/models/company"
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
type employeeDB interface{}
type positionDB interface{}
type departmentDB interface{}
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
