package company_test

import (
	"context"
	"database/sql"
	"fmt"
	"labyrinth/database/postgres"
	"labyrinth/database/postgres/company"
	c "labyrinth/models/company"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

var (
	db          *sql.DB
	testCompany *c.Company
)

func setup() error {
	var connection string = postgres.GetConnection()
	var err error
	db, err = sql.Open("postgres", connection)
	if err != nil {
		return fmt.Errorf("failed to connect to db  during test user: %w", err)
	}

	testCompany = &c.Company{
		ID:          uuid.New(),
		OwnerID:     uuid.New(),
		Name:        "ARASAKA",
		Description: "Create weapons",
		LogoURL:     "https://example.com/araska-logo.png",
		Industry:    "Military Industrial",
		Employees:   25000,
		IsVerified:  true,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		FoundedDate: time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC),
		Address:     "1 Corporate Plaza, Night City",
		Phone:       "+77775553535",
		Email:       "info@arasaka.com",
		TaxNumber:   "US-123456789",
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

func TestCompanyCRUD(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pc := company.NewPostgresCompany()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	defer func() {
		if t.Failed() {
			tx.Rollback()
		}
	}()

	t.Run("CreateCompany", func(t *testing.T) {
		err := pc.CreateCompany(ctx, tx, testCompany)
		if err != nil {
			t.Fatalf("CreateCompany failed: %v", err)
		}
	})

	t.Run("GetCompanyByID", func(t *testing.T) {
		fetchedCompany, err := pc.GetCompanyByID(ctx, tx, testCompany.ID)
		if err != nil {
			t.Fatalf("Failed to get company by ID: %v", err)
		}
		if fetchedCompany.Name != testCompany.Name {
			t.Fatalf("Expected name %q, got %q", testCompany.Name, fetchedCompany.Name)
		}
	})

	t.Run("AddUserToCompany", func(t *testing.T) {
		err := pc.AddUserToCompany(ctx, tx, testCompany.OwnerID, testCompany.ID)
		if err != nil {
			t.Fatalf("Failed to add user to company: %v", err)
		}
	})

	t.Run("GetCompaniesByUser", func(t *testing.T) {
		companies, err := pc.GetCompaniesByUser(ctx, tx, testCompany.OwnerID)
		if err != nil {
			t.Fatalf("Failed to get user companies: %v", err)
		}
		if len(*companies) != 1 {
			t.Fatalf("Expected 1 company, got %d", len(*companies))
		}
		if (*companies)[0].ID != testCompany.ID {
			t.Errorf("Expected company ID %v, got %v", testCompany.ID, (*companies)[0].ID)
		}
	})

	t.Run("UpdateCompany", func(t *testing.T) {
		updatedCompany := *testCompany
		updatedCompany.Name = "ARASAKA UPDATED"
		updatedCompany.Description = "Updated description"

		err := pc.UpdateCompany(ctx, tx, &updatedCompany)
		if err != nil {
			t.Fatalf("Failed to update company: %v", err)
		}

		// Verify update
		fetched, err := pc.GetCompanyByID(ctx, tx, testCompany.ID)
		if err != nil {
			t.Fatalf("Failed to verify update: %v", err)
		}
		if fetched.Name != "ARASAKA UPDATED" {
			t.Errorf("Update failed, expected name %q, got %q", "ARASAKA UPDATED", fetched.Name)
		}
	})

	t.Run("DeleteCompany", func(t *testing.T) {
		err := pc.DeleteCompany(ctx, tx, testCompany.ID)
		if err != nil {
			t.Fatalf("Failed to delete company: %v", err)
		}
	})

	t.Run("DeactivateCompanyUsers", func(t *testing.T) {
		err := pc.DeactivateCompanyUsers(ctx, tx, testCompany.ID)
		if err != nil {
			t.Fatalf("Failed to deactivate company users: %v", err)
		}

		// Verify deactivation
		companies, err := pc.GetCompaniesByUser(ctx, tx, testCompany.OwnerID)
		if err != nil {
			t.Fatalf("Failed to verify deactivation: %v", err)
		}
		if len(*companies) > 0 {
			t.Errorf("Expected 0 companies after deactivation, got %d", len(*companies))
		}
	})

	// Если все тесты прошли успешно, коммитим транзакцию
	if !t.Failed() {
		if err := tx.Rollback(); err != nil {
			t.Errorf("Failed to commit transaction: %v", err)
		}
	}

}
