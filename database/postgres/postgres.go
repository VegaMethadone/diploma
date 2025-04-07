package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"labyrinth/config"
	"labyrinth/models/company"
	"labyrinth/models/department"
	"labyrinth/models/employee"
	"labyrinth/models/position"
	"labyrinth/models/user"

	"github.com/google/uuid"
)

type userDB interface {
	// CreateUser создает нового пользователя (возвращает созданного пользователя с заполненными полями)
	CreateUser(
		ctx context.Context,
		db *sql.DB,
		u *user.User,
	) (*user.User, error)

	// GetUserByCredentials ищет пользователя по логину и проверяет пароль
	GetUserByCredentials(
		ctx context.Context,
		db *sql.DB,
		login,
		password string,
	) (*user.User, error)

	// GetUserByID получает пользователя по ID (без проверки пароля)
	GetUserByID(
		ctx context.Context,
		db *sql.DB,
		id uuid.UUID,
	) (*user.User, error)

	// UpdateUser обновляет данные пользователя
	UpdateUser(
		ctx context.Context,
		db *sql.DB,
		u *user.User,
	) error

	// DeleteUser мягкое удаление (is_active = false)
	DeleteUser(
		ctx context.Context,
		db *sql.DB,
		id uuid.UUID,
	) error
}

type companyDB interface {
	// CreateCompany создает новую компанию и возвращает созданную запись с заполненными полями (ID, даты)
	CreateCompany(
		ctx context.Context,
		db *sql.DB,
		c *company.Company,
	) (*company.Company, error)

	// GetCompanyByID получает компанию по UUID с проверкой существования
	GetCompanyByID(
		ctx context.Context,
		db *sql.DB,
		id uuid.UUID,
	) (*company.Company, error)

	// GetCompanyByOwner получает список компаний по UUID владельца
	GetCompaniesByOwner(
		ctx context.Context,
		db *sql.DB,
		ownerID uuid.UUID,
	) ([]*company.Company, error)

	// UpdateCompany обновляет основные поля компании (название, описание)
	UpdateCompany(
		ctx context.Context,
		db *sql.DB,
		c *company.Company,
	) error

	// DeleteCompany выполняет "мягкое" удаление компании (is_active = false)
	DeleteCompany(
		ctx context.Context,
		db *sql.DB,
		id uuid.UUID,
	) error

	// ForceDeleteCompany полностью удаляет компанию из базы данных
	ForceDeleteCompany(
		ctx context.Context,
		db *sql.DB,
		id uuid.UUID,
	) error

	// AddEmployee добавляет сотрудника в компанию
	// Проверяет что сотрудник не уже добавлен
	AddEmployee(
		ctx context.Context,
		db *sql.DB,
		companyID,
		userID uuid.UUID,
		position string,
	) error

	// RemoveEmployee удаляет сотрудника из компании
	// Возвращает ошибку если сотрудник не найден в компании
	RemoveEmployee(
		ctx context.Context,
		db *sql.DB,
		companyID,
		userID uuid.UUID,
	) error

	// ListEmployees возвращает список всех сотрудников компании с пагинацией
	// Параметр limit = 0 означает отсутствие лимита
	ListEmployees(
		ctx context.Context,
		db *sql.DB,
		companyID uuid.UUID,
		limit,
		offset int,
	) ([]*employee.Employee, error)
}

type employeeDB interface {
	// CreateEmployee создает нового сотрудника и возвращает созданную запись
	// Проверяет уникальность связки user_id + company_id
	CreateEmployee(ctx context.Context, db *sql.DB, emp *employee.Employee) (*employee.Employee, error)

	// GetEmployeeByID получает полную информацию о сотруднике по ID
	// Включая связанные данные пользователя (user) и компании (company)
	GetEmployeeByID(ctx context.Context, db *sql.DB, id uuid.UUID) (*employee.EmployeeWithDetails, error)

	// GetEmployeeByUserCompany получает сотрудника по связке user_id + company_id
	// Используется для проверки "является ли пользователь сотрудником компании"
	GetEmployeeByUserCompany(ctx context.Context, db *sql.DB, userID, companyID uuid.UUID) (*employee.Employee, error)

	// UpdateEmployee обновляет основные данные сотрудника
	// (должность, уровень доступа, статус активности)
	UpdateEmployee(ctx context.Context, db *sql.DB, emp *employee.Employee) error

	// UpdateEmployeeStatus изменяет статус активности сотрудника (is_active)
	// Для архивных записей сохраняет дату деактивации (deactivated_at)
	UpdateEmployeeStatus(ctx context.Context, db *sql.DB, id uuid.UUID, isActive bool) error

	// ListCompanyEmployees возвращает всех сотрудников компании
	// Поддерживает фильтрацию по:
	// - статусу (active/inactive)
	// - должности
	// - уровню доступа
	ListCompanyEmployees(
		ctx context.Context,
		db *sql.DB,
		companyID uuid.UUID,
		filter EmployeeFilter,
	) ([]*employee.EmployeeWithUserDetails, error)

	// ListUserCompanies возвращает все компании, в которых состоит пользователь
	// С указанием его должности в каждой компании
	ListUserCompanies(ctx context.Context, db *sql.DB, userID uuid.UUID) ([]*employee.UserCompany, error)

	// DeleteEmployee полностью удаляет запись сотрудника
	// Использовать только при ошибочном добавлении!
	// Для обычного "увольнения" использовать UpdateEmployeeStatus(false)
	DeleteEmployee(ctx context.Context, db *sql.DB, id uuid.UUID) error
}

/*
	Сомнительно
*/
// EmployeeFilter - фильтр для списка сотрудников
type EmployeeFilter struct {
	IsActive    *bool       // Фильтр по статусу (nil - все)
	PositionIDs []uuid.UUID // Фильтр по должностям
	AccessLevel *int        // Фильтр по уровню доступа
	SearchQuery string      // Поиск по имени/фамилии/email
}
type positionDB interface {
	// CreatePosition создает новую должность в указанном департаменте
	// Проверяет уникальность названия в рамках департамента
	CreatePosition(ctx context.Context, db *sql.DB, position *position.Position) (*position.Position, error)

	// GetPositionByID получает должность по ID с проверкой существования
	// Возвращает ошибку sql.ErrNoRows если должность не найдена
	GetPositionByID(ctx context.Context, db *sql.DB, id uuid.UUID) (*position.Position, error)

	// GetDepartmentPositions возвращает все должности в департаменте
	// Сортировка по уровню доступа (от высшего к низшему)
	GetDepartmentPositions(ctx context.Context, db *sql.DB, departmentID uuid.UUID) ([]*position.Position, error)

	// UpdatePosition обновляет данные должности:
	// - название
	// - описание
	// - уровень доступа
	// Не позволяет изменять department_id (для этого нужно удалить/создать)
	UpdatePosition(ctx context.Context, db *sql.DB, position *position.Position) error

	// UpdatePositionAccess изменяет уровень доступа должности
	// Автоматически обновляет права всех сотрудников на этой должности
	UpdatePositionAccess(ctx context.Context, db *sql.DB, positionID uuid.UUID, newAccessLevel int) error

	// DeletePosition удаляет должность если на ней нет сотрудников
	// Возвращает ошибку если:
	// - должность не найдена
	// - на должности есть сотрудники
	DeletePosition(ctx context.Context, db *sql.DB, id uuid.UUID) error

	// ForceDeletePosition принудительно удаляет должность и переносит
	// всех сотрудников на указанную целевую должность
	ForceDeletePosition(
		ctx context.Context,
		db *sql.DB,
		positionID uuid.UUID,
		moveToPositionID uuid.UUID,
	) error

	// GetCompanyPositions возвращает все должности компании
	// (агрегирует позиции из всех департаментов)
	// Поддерживает фильтрацию по уровню доступа
	GetCompanyPositions(
		ctx context.Context,
		db *sql.DB,
		companyID uuid.UUID,
		filter PositionFilter,
	) ([]*position.PositionWithDepartment, error)
}

// PositionFilter - фильтр для списка должностей
type PositionFilter struct {
	MinAccessLevel *int       // Минимальный уровень доступа
	MaxAccessLevel *int       // Максимальный уровень доступ
	DepartmentID   *uuid.UUID // Фильтр по департаменту
	OnlyVacant     bool       // Только должности без сотрудников
}
type departmentDB interface {
	// CreateDepartment создает новый департамент в компании
	// Проверяет уникальность названия в рамках компании
	CreateDepartment(
		ctx context.Context,
		db *sql.DB,
		dep *department.Department,
	) (*department.Department, error)

	// GetDepartmentByID получает департамент по ID с полной информацией:
	// - базовые данные департамента
	// - информация о родительском департаменте (если есть)
	// - количество сотрудников
	GetDepartmentByID(
		ctx context.Context,
		db *sql.DB,
		id uuid.UUID,
	) (*department.DepartmentWithDetails, error)

	// GetCompanyDepartments возвращает иерархию департаментов компании
	// Параметр flat=false возвращает древовидную структуру
	// Параметр flat=true возвращает плоский список
	GetCompanyDepartments(
		ctx context.Context,
		db *sql.DB,
		companyID uuid.UUID,
		flat bool,
	) ([]*department.DepartmentNode, error)

	// UpdateDepartment обновляет основные данные департамента:
	// - название
	// - описание
	// - руководитель (owner_employee_id)
	// Не позволяет изменять company_id
	UpdateDepartment(
		ctx context.Context,
		db *sql.DB,
		dep *department.Department,
	) error

	// MoveDepartment перемещает департамент в другую родительскую ветку
	// Проверяет что:
	// - новый родитель существует
	// - не создается циклическая ссылка
	MoveDepartment(
		ctx context.Context,
		db *sql.DB,
		departmentID uuid.UUID,
		newParentID *uuid.UUID, // nil для корневого уровня
	) error

	// DeleteDepartment удаляет департамент если:
	// - нет дочерних департаментов
	// - нет сотрудников
	// Для принудительного удаления использовать ForceDeleteDepartment
	DeleteDepartment(ctx context.Context, db *sql.DB, id uuid.UUID) error

	// ForceDeleteDepartment принудительно удаляет департамент:
	// - переносит дочерние департаменты на уровень выше
	// - перемещает сотрудников в указанный департамент
	ForceDeleteDepartment(
		ctx context.Context,
		db *sql.DB,
		departmentID uuid.UUID,
		moveEmployeesTo *uuid.UUID, // куда перемещать сотрудников
	) error

	// GetDepartmentEmployees возвращает сотрудников департамента
	// с возможностью фильтрации по должностям
	GetDepartmentEmployees(
		ctx context.Context,
		db *sql.DB,
		departmentID uuid.UUID,
		filter DepartmentEmployeeFilter,
	) ([]*department.EmployeeWithPosition, error)

	// GetEmployeeDepartments возвращает все департаменты сотрудника
	// в указанной компании (с учетом множественного прикрепления)
	GetEmployeeDepartments(
		ctx context.Context,
		db *sql.DB,
		companyID uuid.UUID,
		employeeID uuid.UUID,
	) ([]*department.DepartmentWithPosition, error)
}

type departmentEmployeeDB interface {
	// AssignEmployeeToDepartment добавляет сотрудника в департамент
	// Проверяет что сотрудник уже не прикреплен к этому департаменту
	AssignEmployeeToDepartment(
		ctx context.Context,
		db *sql.DB,
		departmentID uuid.UUID,
		employeeID uuid.UUID,
	) error

	// RemoveEmployeeFromDepartment удаляет прикрепление сотрудника к департаменту
	// Проверяет существование связи перед удалением
	RemoveEmployeeFromDepartment(
		ctx context.Context,
		db *sql.DB,
		departmentID uuid.UUID,
		employeeID uuid.UUID,
	) error

	// GetEmployeeDepartments возвращает список департаментов сотрудника
	// с указанием основной должности в каждом
	GetEmployeeDepartments(
		ctx context.Context,
		db *sql.DB,
		employeeID uuid.UUID,
	) ([]*DepartmentAssignment, error)

	// GetDepartmentEmployees возвращает список сотрудников департамента
	// с возможностью фильтрации по статусу активности
	GetDepartmentEmployees(
		ctx context.Context,
		db *sql.DB,
		departmentID uuid.UUID,
		onlyActive bool,
	) ([]*EmployeeAssignment, error)

	// UpdateEmployeePosition изменяет основную должность сотрудника в департаменте
	// Проверяет что должность существует в этом департаменте
	UpdateEmployeePosition(
		ctx context.Context,
		db *sql.DB,
		departmentID uuid.UUID,
		employeeID uuid.UUID,
		newPositionID uuid.UUID,
	) error

	// IsEmployeeInDepartment проверяет принадлежность сотрудника к департаменту
	// Возвращает false если связи не существует
	IsEmployeeInDepartment(
		ctx context.Context,
		db *sql.DB,
		departmentID uuid.UUID,
		employeeID uuid.UUID,
	) (bool, error)

	// GetPrimaryDepartment возвращает основной департамент сотрудника
	// (помеченный как is_primary) или nil если не назначен
	GetPrimaryDepartment(
		ctx context.Context,
		db *sql.DB,
		employeeID uuid.UUID,
	) (*uuid.UUID, error)

	// SetPrimaryDepartment устанавливает основной департамент для сотрудника
	// Автоматически снимает пометку с предыдущего основного департамента
	SetPrimaryDepartment(
		ctx context.Context,
		db *sql.DB,
		departmentID uuid.UUID,
		employeeID uuid.UUID,
	) error
}

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
