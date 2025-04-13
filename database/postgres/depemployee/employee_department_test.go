package depemployee_test

import (
	"context"
	"database/sql"
	"fmt"
	"labyrinth/database/postgres"
	"labyrinth/database/postgres/depemployee"
	d "labyrinth/models/depemployee"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

var (
	db              *sql.DB
	testDepEmployee *d.DepartmentEmployee
)

func setup() error {
	var connection string = postgres.GetConnection()
	var err error
	db, err = sql.Open("postgres", connection)
	if err != nil {
		return fmt.Errorf("failed to connect to db  during test user: %w", err)
	}
	testDepEmployee = &d.DepartmentEmployee{
		ID:           uuid.New(),
		EmployeeID:   uuid.New(),
		DepartmentID: uuid.New(),
		PositionID:   uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsActive:     true,
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

func TestDepartmentEmployeeCRUD(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pde := depemployee.NewPostgresEmployeeDepartment()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v\n", err)
	}
	defer func() {
		if t.Failed() {
			tx.Rollback()
		}
	}()

	t.Run("CreateEmployeeDepartment", func(t *testing.T) {
		err = pde.CreateEmployeeDepartment(ctx, tx, testDepEmployee)
		if err != nil {
			t.Fatalf("CreateEmployeeDepartment failed: %v\n", err)
		}
	})

	t.Run("ExistsEmployeeDepartment", func(t *testing.T) {
		exists, err := pde.ExistsEmployeeDepartment(ctx, tx, testDepEmployee.EmployeeID, testDepEmployee.DepartmentID)
		if err != nil {
			t.Fatalf("ExistsEmployeeDepartment failed: %v\n", err)
		}

		if !exists {
			t.Errorf("Expected true, got false\n")
		}
	})

	t.Run("UpdateEmployeeDepartment", func(t *testing.T) {
		updatedDepEmployee := *testDepEmployee
		updatedDepEmployee.PositionID = uuid.New()
		updatedDepEmployee.UpdatedAt = time.Now()

		err = pde.UpdateEmployeeDepartment(ctx, tx, &updatedDepEmployee)
		if err != nil {
			t.Fatalf("UpdateEmployeeDepartment failed: %v\n", err)
		}
	})

	t.Run("GetEmployeeDepartmentByEmployeeId", func(t *testing.T) {
		fetchedDepemployee, err := pde.GetEmployeeDepartmentByEmployeeId(ctx, tx, testDepEmployee.EmployeeID, testDepEmployee.DepartmentID)
		if err != nil {
			t.Fatalf("GetEmployeeDepartmentByEmployeeId failed: %v\n", err)
		}

		if fetchedDepemployee.ID != testDepEmployee.ID {
			t.Errorf("Expected %s, got %s\n", testDepEmployee.ID, fetchedDepemployee.ID)
		}
	})

	t.Run("GetEmployeesDepartmentByDepartmentId", func(t *testing.T) {
		fetchedEmployees, err := pde.GetEmployeesDepartmentByDepartmentId(ctx, tx, testDepEmployee.DepartmentID)
		if err != nil {
			t.Fatalf("GetEmployeesDepartmentByDepartmentId failed: %v\n", err)
		}

		if len(fetchedEmployees) != 1 {
			t.Errorf("Expected 1, got %d\n", len(fetchedEmployees))
		}
	})

	t.Run("DeleteEmployeeDepartment", func(t *testing.T) {
		err = pde.DeleteEmployeeDepartment(ctx, tx, testDepEmployee.ID)
		if err != nil {
			t.Fatalf("DeleteEmployeeDepartment failed: %v\n", err)
		}
	})

	if !t.Failed() {
		if err = tx.Rollback(); err != nil {
			t.Errorf("Failled to rollback transaction: %v\n", err)
		}
	}
}
