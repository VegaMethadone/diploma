package position_test

import (
	"context"
	"database/sql"
	"fmt"
	"labyrinth/database/postgres"
	"labyrinth/database/postgres/position"
	p "labyrinth/models/position"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

var (
	db           *sql.DB
	testPosition *p.Position
)

func setup() error {
	var connection string = postgres.GetConnection()
	var err error
	db, err = sql.Open("postgres", connection)
	if err != nil {
		return fmt.Errorf("failed to connect to db  during test position: %w", err)
	}
	testPosition = &p.Position{
		ID:        uuid.New(),
		CompanyID: uuid.New(),
		Lvl:       1,
		Name:      "Ping-Pong",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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

func TestPositionCRUD(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pp := position.NewPostgresPosition()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	defer func() {
		if t.Failed() {
			tx.Rollback()
		}
	}()

	t.Run("CreatePosition", func(t *testing.T) {
		err := pp.CreatePosition(ctx, tx, testPosition)
		if err != nil {
			t.Fatalf("CreatePosition failed: %v", err)
		}
	})

	t.Run("GetPositionById", func(t *testing.T) {
		fetchedPosition, err := pp.GetPositionById(ctx, tx, testPosition.ID)
		if err != nil {
			t.Fatalf("GetPositionById failed: %v", err)
		}

		if fetchedPosition.Name != testPosition.Name {
			t.Errorf("Expected name %q, got %q", testPosition.Name, fetchedPosition.Name)
		}
	})

	t.Run("GetPositionsByCompanyId", func(t *testing.T) {
		fetchedPositions, err := pp.GetPositionsByCompanyId(ctx, tx, testPosition.CompanyID)
		if err != nil {
			t.Fatalf("GetPositionsByCompanyId failed: %v", err)
		}

		for _, value := range *fetchedPositions {
			if value.CompanyID != testPosition.CompanyID {
				t.Errorf("Expected companyId %q, got %q", testPosition.CompanyID, value.CompanyID)
			}
		}
	})

	t.Run("UpdatePosition", func(t *testing.T) {
		updatedPosition := *testPosition
		updatedPosition.Name = "UpdatedPosition"

		err := pp.UpdatePosition(ctx, tx, &updatedPosition)
		if err != nil {
			t.Fatalf("UpdatePosition failed: %v", err)
		}
	})

	t.Run("DeletePosition", func(t *testing.T) {
		err := pp.DeletePosition(ctx, tx, testPosition.ID)
		if err != nil {
			t.Fatalf("DeletePosition failed: %v", err)
		}
	})

	if !t.Failed() {
		if err := tx.Rollback(); err != nil {
			t.Errorf("Failed to rollback transaction: %v", err)
		}
	}
}
