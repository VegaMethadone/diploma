package employee_test

import (
	"context"
	"database/sql"
	"fmt"
	"labyrinth/database/postgres"
	"labyrinth/database/postgres/employee"
	e "labyrinth/models/employee"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

var (
	db           *sql.DB
	testEmployee *e.Employee
)

func setup() error {
	var connection string = postgres.GetConnection()
	var err error
	db, err = sql.Open("postgres", connection)
	if err != nil {
		return fmt.Errorf("failed to connect to db  during test user: %w", err)
	}
	testEmployee = &e.Employee{
		ID:             uuid.New(),
		UserID:         uuid.New(),
		CompanyID:      uuid.New(),
		PositionID:     uuid.New(),
		IsActive:       true,
		IsOnline:       false,
		LastActivityAt: time.Now(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return nil
}

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		fmt.Printf("Test setup failed: %v\n", err)
		os.Exit(1)
	}
	code := m.Run()

	teardown()
	os.Exit(code)

}
func teardown() {
	if db != nil {
		db.Close()
	}
}

func TestEmployeeCRUD(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pe := employee.NewPostgresEmployee()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v\n", err)
	}
	defer func() {
		if t.Failed() {
			tx.Rollback()
		}
	}()

	t.Run("CreateEmployee", func(t *testing.T) {
		err = pe.CreateEmployee(ctx, tx, testEmployee)
		if err != nil {
			t.Fatalf("CreateEmployee failed: %v\n", err)
		}
	})

	t.Run("CountEmployees", func(t *testing.T) {
		count, err := pe.CountEmployees(ctx, tx, testEmployee.CompanyID)
		if err != nil {
			t.Fatalf("CountEmployees failed: %v\n", err)
		}

		if count != 1 {
			t.Errorf("Expected employee count 1, got %d\n", count)
		}
	})

	t.Run("ExistsEmployee", func(t *testing.T) {
		exists, err := pe.ExistsEmployee(ctx, tx, testEmployee.ID)
		if err != nil {
			t.Fatalf("ExistsEmployee failed: %v\n", err)
		}

		if !exists {
			t.Errorf("Expected true, got false\n")
		}
	})

	t.Run("GetEmployeesByCompanyId", func(t *testing.T) {
		fetchedEmployees, err := pe.GetEmployeesByCompanyId(ctx, tx, testEmployee.CompanyID)
		if err != nil {
			t.Fatalf("GetEmployeesByCompanyId failed: %v\n", err)
		}

		if len(*fetchedEmployees) != 1 {
			t.Errorf("Expected size of arr 1, got %d\n", len(*fetchedEmployees))
		}
	})

	t.Run("UpdateEmployee", func(t *testing.T) {
		updatedEmployee := *testEmployee
		updatedEmployee.UpdatedAt = time.Now()
		updatedEmployee.IsActive = false

		err = pe.UpdateEmployee(ctx, tx, &updatedEmployee)
		if err != nil {
			t.Fatalf("UpdateEmployee failed: %v\n", err)
		}
	})

	t.Run("GetEmployeeByUserId", func(t *testing.T) {
		fetchedEmployee, err := pe.GetEmployeeByUserId(ctx, tx, testEmployee.UserID, testEmployee.CompanyID)
		if err != nil {
			t.Fatalf("GetEmployeeByUserId failed: %v\n", err)
		}
		if fetchedEmployee.IsActive {
			t.Errorf("Expected false, got true\n")
		}
	})

	if !t.Failed() {
		if err = tx.Rollback(); err != nil {
			t.Errorf("Failed to rollback transaction: %v\n", err)
		}
	}
}
