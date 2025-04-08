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
		EmailVerified:    false,
		Phone:            "+77775553535",
		PhoneVerified:    false,
		FirstName:        "ivan",
		LastName:         "ivanovich",
		Bio:              "NoBio",
		TelegramUsername: "@ivanovich",
		AvatarURL:        "/c/users/ivan/gay.jpg",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		LastLoginAt:      time.Now(),
		IsActive:         false,
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
		err := user.CreateUser(ctx, tx, testUser)
		if err != nil {
			t.Fatalf("CreateUser failed: %v", err)
		}
	})

	t.Run("GetUser", func(t *testing.T) {
		fetchedUser, err := user.GetUserByID(ctx, tx, testUser.ID)
		if err != nil {
			t.Fatalf("GetUserByID failed: %v", err)
		}

		if fetchedUser.Login != testUser.Login {
			t.Errorf("Expected login %q, got %q", testUser.Login, fetchedUser.Login)
		}
	})

	t.Run("GetUserByCredentials", func(t *testing.T) {
		fetchedUser, err := user.GetUserByCredentials(ctx, tx, testUser.Login, testUser.PasswordHash)
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

		err := user.UpdateUser(ctx, tx, &updatedUser)
		if err != nil {
			t.Fatalf("UpdateUser failed: %v", err)
		}
	})

	t.Run("DeleteUser", func(t *testing.T) {
		err := user.DeleteUser(ctx, tx, testUser.ID)
		if err != nil {
			t.Fatalf("DeleteUser failed: %v", err)
		}
	})

	// Если все тесты прошли успешно, коммитим транзакцию
	if !t.Failed() {
		if err := tx.Commit(); err != nil {
			t.Errorf("Failed to commit transaction: %v", err)
		}
	}
}
