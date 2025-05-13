package userlogic_test

import (
	"fmt"
	"labyrinth/logger"
	authlogic "labyrinth/logic/authLogic"
	userlogic "labyrinth/logic/userLogic"
	"labyrinth/models/user"
	"os"
	"testing"
)

var (
	auth     authlogic.Auth      = authlogic.NewAuth()
	usr      userlogic.Userlogic = userlogic.NewUserlogic()
	testUser *user.User          = user.NewUser(
		"user_test1@gmail.com",
		"1234567891111",
		"+75555553566",
	)

	fetchedUser *user.User
)

func setup() error {
	err := auth.Register(testUser.Email, testUser.PasswordHash, testUser.Phone)
	if err != nil {
		return err
	}
	fetchedUser, err = auth.Login(testUser.Email, testUser.PasswordHash)
	if err != nil {
		return err
	}

	return nil
}

func TestMain(m *testing.M) {
	logger.InitFileLogger("user_test.logs")
	if err := setup(); err != nil {
		fmt.Printf("Test setup failed: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	os.Exit(code)
}

func TestUser(t *testing.T) {

	t.Run("UpdateUserProfile", func(t *testing.T) {
		fetchedUser.Bio = "DEAD SPACE"
		err := usr.UpdateUserProfile(fetchedUser)
		if err != nil {
			t.Fatalf("Failed UpdateUserProfile: %v", err)
		}
	})

	t.Run("GetUserProfile", func(t *testing.T) {
		gotUser, err := usr.GetUserProfile(fetchedUser.ID)
		if err != nil {
			t.Fatalf("Failed GetUserProfile: %v", err)
		}

		if gotUser.Bio != "DEAD SPACE" {
			t.Errorf("Expected DEAD SPACE, got %s\n", gotUser.Bio)
		}
	})
}
