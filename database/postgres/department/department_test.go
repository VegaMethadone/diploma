package department_test

import (
	"context"
	"database/sql"
	"fmt"
	"labyrinth/database/postgres"
	"labyrinth/database/postgres/department"
	d "labyrinth/models/department"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

var (
	db             *sql.DB
	testDepartment *d.Department
)

func setup() error {
	var connection string = postgres.GetConnection()
	var err error
	db, err = sql.Open("postgres", connection)
	if err != nil {
		return fmt.Errorf("failed to connect to db  during test user: %w", err)
	}
	parentId := uuid.New()
	testDepartment = &d.Department{
		ID:          uuid.New(),
		CompanyID:   parentId,
		Name:        "ARASAKA DEPARTMENT",
		Description: "MILITARY WEAPON",
		AvatarURL:   "C/images/dep.jpg",
		ParentID:    parentId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsActive:    true,
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

func TestDepartmentCRUD(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pd := department.NewPostgresDepartment()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v\n", err)
	}
	defer func() {
		if t.Failed() {
			tx.Rollback()
		}
	}()

	t.Run("CreateDepartment", func(t *testing.T) {
		err = pd.CreateDepartment(ctx, tx, testDepartment)
		if err != nil {
			t.Fatalf("CreateDepartment failed: %v\n", err)
		}
	})

	t.Run("UpdateDepartment", func(t *testing.T) {
		updatedDepartment := *testDepartment
		updatedDepartment.Name = "ARASAKA DEPARTMENT UPDATED"
		updatedDepartment.UpdatedAt = time.Now()

		err = pd.UpdateDepartment(ctx, tx, &updatedDepartment)
		if err != nil {
			t.Fatalf("UpdateDepartment failed:  %v\n", err)
		}
	})

	t.Run("GetDepartmentById", func(t *testing.T) {
		fetchedDepartment, err := pd.GetDepartmentById(ctx, tx, testDepartment.ID)
		if err != nil {
			t.Fatalf("GetDepartmentById failed %v\n", err)
		}

		if fetchedDepartment.Name != "ARASAKA DEPARTMENT UPDATED" {
			t.Errorf("Expected ARASAKA DEPARTMENT UPDATED, got  %s\n", fetchedDepartment.Name)
		}
	})

	t.Run("GetDepartmentsByParentId", func(t *testing.T) {
		fetchedDepartments, err := pd.GetDepartmentsByParentId(ctx, tx, testDepartment.ParentID)
		if err != nil {
			t.Fatalf("GetDepartmentsByParentId failed: %v\n", err)
		}

		if len(fetchedDepartments) != 1 {
			t.Errorf("Expected size 1, got: %d\n", len(fetchedDepartments))
		}
	})

	t.Run("DeleteDepartment", func(t *testing.T) {
		err = pd.DeleteDepartment(ctx, tx, testDepartment.ID)
		if err != nil {
			t.Fatalf("DeleteDepartment failed: %v\n", err)
		}
	})

	if !t.Failed() {
		if err = tx.Rollback(); err != nil {
			t.Errorf("Failed to rollback transaction: %v\n", err)
		}
	}
}
