package positionlogic_test

import (
	"fmt"
	"labyrinth/logger"
	authlogic "labyrinth/logic/authLogic"
	companylogic "labyrinth/logic/companyLogic"
	positionlogic "labyrinth/logic/positionLogic"
	"labyrinth/models/company"
	"labyrinth/models/position"
	"labyrinth/models/user"
	"os"
	"testing"

	"github.com/google/uuid"
)

var (
	auth           authlogic.Auth              = authlogic.NewAuth()
	comp           companylogic.CompanyLogic   = companylogic.NewCompanyLogic()
	pos            positionlogic.PositionLogic = positionlogic.NewPositionLogic()
	fetchedCompany *company.Company
	fetched1User   *user.User
	testUser       *user.User = user.NewUser(
		"position_test1@gmail.com",
		"123456789",
		"+77855553611",
	)
)

func setup() error {
	err := auth.Register(testUser.Login, testUser.PasswordHash, testUser.Phone)
	if err != nil {
		return err
	}

	fetched1User, err = auth.Login(testUser.Login, testUser.PasswordHash)
	if err != nil {
		return err
	}

	companyId, err := comp.NewCompany(fetched1User.ID, "myDepartment", "myDepartment")
	if err != nil {
		return err
	}

	fetchedCompany, err = comp.GetCompany(fetched1User.ID, companyId)
	if err != nil {
		return err
	}

	return nil
}

func TestMain(m *testing.M) {
	logger.InitFileLogger("position_test.logs")
	if err := setup(); err != nil {
		fmt.Printf("Test setup failed: %v\n", err)
		os.Exit(1)
	}
	code := m.Run()

	os.Exit(code)
}

func TestPosition(t *testing.T) {
	positionId := uuid.Nil
	var err error
	t.Run("NewPosition", func(t *testing.T) {
		positionId, err = pos.NewPosition(fetched1User.ID, fetchedCompany.ID, 1, "THE GOD")
		if err != nil {
			t.Fatalf("Failed  NewPosition: %v", err)
		}
	})

	var updatePosition *position.Position
	t.Run("GetAllPositions", func(t *testing.T) {
		positions, err := pos.GetAllPositions(fetched1User.ID, fetchedCompany.ID)
		if err != nil {
			t.Fatalf("Failed GetAllPositions: %v", err)
		}
		for _, position := range *positions {
			if position.ID == positionId {
				updatePosition = &position
			}
		}
		if updatePosition == nil {
			t.Errorf("Expected position with id: %s", positionId.String())
		}
	})

	t.Run("UpdatePosition", func(t *testing.T) {
		updatePosition.Name = "THE GOD UPDATED"
		err = pos.UpdatePosition(fetched1User.ID, fetchedCompany.ID, updatePosition)
		if err != nil {
			t.Fatalf("Failed UpdatePosition: %v", err)
		}
	})
}
