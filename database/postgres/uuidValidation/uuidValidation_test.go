package uuidvalidation_test

import (
	"context"
	"database/sql"
	"fmt"
	"labyrinth/database/postgres"
	uuidvalidation "labyrinth/database/postgres/uuidValidation"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

var (
	db *sql.DB
)

func setup() error {
	var connection string = postgres.GetConnection()
	var err error
	db, err = sql.Open("postgres", connection)
	if err != nil {
		return fmt.Errorf("failed to connect to db  during test user: %w", err)
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

func TestUUIDValidation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	uv := uuidvalidation.NewDBUuidValidation()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	defer func() {
		if t.Failed() {
			tx.Rollback()
		}
	}()

	t.Run("CheckAndReserveUUID", func(t *testing.T) {
		generatedUUID, err := uv.CheckAndReserveUUID(ctx, tx)
		if err != nil {
			t.Fatalf("CheckAndReserveUUID failed: %v", err)
		}

		err = uuid.Validate(generatedUUID.String())
		if err != nil {
			t.Errorf("Failed to validate generated UUID: %v", err)
		} else {
			fmt.Printf("Generated UUID: %s\n", generatedUUID.String())
		}
	})
}
