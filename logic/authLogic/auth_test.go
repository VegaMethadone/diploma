package authlogic_test

import (
	"labyrinth/logger"
	authlogic "labyrinth/logic/authLogic"
	"labyrinth/models/user"
	"os"
	"testing"
)

var (
	auth     authlogic.Auth = authlogic.NewAuth()
	testUser *user.User     = user.NewUser(
		"ivannnn@gmail.com",
		"123456789",
		"+75555553535",
	)
)

func TestMain(m *testing.M) {
	logger.InitFileLogger("auth_test.logs")
	code := m.Run()
	os.Exit(code)
}

func TestAuth(t *testing.T) {
	t.Run("Register", func(t *testing.T) {
		err := auth.Register(testUser.Email, testUser.PasswordHash, testUser.Phone)
		if err != nil {
			t.Fatalf("Failed to register user: %v", err)
		}
	})

	t.Run("Login", func(t *testing.T) {
		fetchedUser, err := auth.Login(testUser.Email, testUser.PasswordHash)
		if err != nil {
			t.Fatalf("Failed to login user: %v", err)
		}

		if fetchedUser.Email != testUser.Email {
			t.Errorf("Expected %s, got %s\n", testUser.Email, fetchedUser.Email)
		}
	})
}
