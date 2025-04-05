package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/database/postgres"
	"labyrinth/entity/employee"

	"github.com/google/uuid"
)

func GetEmployee(userId, companyId uuid.UUID) (*employee.Employee, error) {
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	findQuery := `
		SELECT id, position
		FROM employee_company  
		WHERE user_id = $1 and company_id = $2
	`
	var (
		employeeId uuid.UUID
		positionId uuid.UUID
	)
	err = db.QueryRow(findQuery, userId, companyId).Scan(&employeeId, &positionId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to check employee: %w", err)
	}

	return &employee.Employee{
		Id:         employeeId,
		UserId:     userId,
		CompanyId:  companyId,
		PositionId: positionId,
	}, nil
}
