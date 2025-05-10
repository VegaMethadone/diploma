package depposition_test

import (
	"context"
	"database/sql"
	"fmt"
	"labyrinth/database/postgres"
	"labyrinth/database/postgres/depposition"
	d "labyrinth/models/depposition"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

var (
	db              *sql.DB
	testDepPosition *d.DepPosition
)

func setup() error {
	var connection string = postgres.GetConnection()
	var err error
	db, err = sql.Open("postgres", connection)
	if err != nil {
		return fmt.Errorf("failed to connect to db  during test user: %w", err)
	}
	testDepPosition = &d.DepPosition{
		Id:           uuid.New(),
		DepartmentId: uuid.New(),
		Level:        2,
		Name:         "THE GOD",
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

func TestDepPositionCRUD(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dp := depposition.NewPostgresDepPosition()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v\n", err)
	}
	defer func() {
		if t.Failed() {
			tx.Rollback()
		}
	}()

	t.Run("CreateDepartmentPosition", func(t *testing.T) {
		err = dp.CreateDepartmentPosition(ctx, tx, testDepPosition)
		if err != nil {
			t.Fatalf("CreateDepartmentPosition failed:  %v\n", err)
		}
	})

	t.Run("ExistsPosition", func(t *testing.T) {
		exists, err := dp.ExistsPosition(ctx, tx, testDepPosition.Id)
		if err != nil {
			t.Fatalf("ExistsPosition failed:  %v\n", err)
		}

		if !exists {
			t.Errorf("Expected true, got  false\n")
		}
	})

	t.Run("UpdateDepartmentPosition", func(t *testing.T) {
		updatedPosition := *testDepPosition
		updatedPosition.Name = "UPDATED THE GOD"

		err = dp.UpdateDepartmentPosition(ctx, tx, &updatedPosition)
		if err != nil {
			t.Fatalf("UpdateDepartmentPosition failed: %v\n", err)
		}
	})

	t.Run("GetDepartmentPositionById", func(t *testing.T) {
		fetchedDepPosition, err := dp.GetDepartmentPositionById(ctx, tx, testDepPosition.Id)
		if err != nil {
			t.Fatalf("GetDepartmentPositionById failed: %v\n", err)
		}

		if fetchedDepPosition.Name != "UPDATED THE GOD" {
			t.Errorf("Expected UPDATED THE GOD, got %s\n", fetchedDepPosition.Name)
		}
	})

	t.Run("GetDepartmentPositionsByDepartmentId", func(t *testing.T) {
		fetchedDepPositions, err := dp.GetDepartmentPositionsByDepartmentId(ctx, tx, testDepPosition.DepartmentId)
		if err != nil {
			t.Fatalf("GetDepartmentPositionsByDepartmentId failed: %v\n", err)
		}

		if len(*fetchedDepPositions) != 1 {
			t.Errorf("Expected 1, got %d\n", len(*fetchedDepPositions))
		}
	})

	t.Run("DeleteDepartmentPosition", func(t *testing.T) {
		err = dp.DeleteDepartmentPosition(ctx, tx, testDepPosition.Id)
		if err != nil {
			t.Fatalf("DeleteDepartmentPosition  failed: %v\n", err)
		}
	})

	if !t.Failed() {
		if err = tx.Rollback(); err != nil {
			t.Errorf("Failed to rollback transaction: %v\n", err)
		}
	}
}
