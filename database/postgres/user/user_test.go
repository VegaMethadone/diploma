package user_test

import (
	"context"
	"database/sql"
	"fmt"
	"labyrinth/database/postgres"
	"labyrinth/database/postgres/user"
	u "labyrinth/models/user"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

var (
	db       *sql.DB
	testUser *u.User
)

func setup() error {
	var connection string = postgres.GetConnection()
	var err error
	db, err = sql.Open("postgres", connection)
	if err != nil {
		return fmt.Errorf("failed to connect to db  during test user: %w", err)
	}
	testUser = &u.User{
		ID:               uuid.New(),
		Login:            "vegaomega",
		PasswordHash:     "qweqwe123",
		Email:            "vegaomega@gmail.com",
		EmailVerified:    true,
		Phone:            "+77775553535",
		PhoneVerified:    true,
		FirstName:        "ivan",
		LastName:         "ivanovich",
		Bio:              "NoBio",
		TelegramUsername: "@ivanovich",
		AvatarURL:        "/c/users/ivan/gay.jpg",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		LastLoginAt:      time.Now(),
		IsActive:         true,
		IsStaff:          false,
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
func TestUserCRUD(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pu := user.NewPostgresUser()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	defer func() {
		if t.Failed() {
			tx.Rollback()
		}
	}()

	t.Run("CreateUser", func(t *testing.T) {
		err := pu.CreateUser(ctx, tx, testUser)
		if err != nil {
			t.Fatalf("CreateUser failed: %v", err)
		}
	})

	t.Run("GetUser", func(t *testing.T) {
		fetchedUser, err := pu.GetUserByID(ctx, tx, testUser.ID)
		if err != nil {
			t.Fatalf("GetUserByID failed: %v", err)
		}

		if fetchedUser.Login != testUser.Login {
			t.Errorf("Expected login %q, got %q", testUser.Login, fetchedUser.Login)
		}
	})

	t.Run("GetUserByCredentials", func(t *testing.T) {
		fetchedUser, err := pu.GetUserByCredentials(ctx, tx, testUser.Login, testUser.PasswordHash)
		if err != nil {
			t.Fatalf("GetUserByID failed: %v", err)
		}
		if fetchedUser.Login != testUser.Login {
			t.Errorf("Expected login %q, got %q", testUser.Login, fetchedUser.Login)
		}
	})

	t.Run("UpdateUser", func(t *testing.T) {
		updatedUser := *testUser
		updatedUser.FirstName = "UpdatedName"

		err := pu.UpdateUser(ctx, tx, &updatedUser)
		if err != nil {
			t.Fatalf("UpdateUser failed: %v", err)
		}
	})

	t.Run("CheckPhone", func(t *testing.T) {
		exists, err := pu.CheckPhone(ctx, tx, "+77775553535")
		if err != nil {
			t.Fatalf("CheckPhone failed: %v", err)
		}

		if !exists {
			t.Errorf("Expected phone exists = true, got phone exists = false")
		}
	})

	t.Run("DeleteUser", func(t *testing.T) {
		err := pu.DeleteUser(ctx, tx, testUser.ID)
		if err != nil {
			t.Fatalf("DeleteUser failed: %v", err)
		}
	})

	if !t.Failed() {
		if err := tx.Rollback(); err != nil {
			t.Errorf("Failed to rollback transaction: %v", err)
		}
	}
}
