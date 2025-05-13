package companylogic_test

import (
	"labyrinth/logger"
	authlogic "labyrinth/logic/authLogic"
	companylogic "labyrinth/logic/companyLogic"
	"labyrinth/models/company"
	"labyrinth/models/user"
	"os"
	"testing"

	"github.com/google/uuid"
)

var (
	comp           companylogic.CompanyLogic = companylogic.NewCompanyLogic()
	auth           authlogic.Auth            = authlogic.NewAuth()
	name           string                    = "myCompany"
	description    string                    = "myCompany"
	userId         uuid.UUID
	companyId      uuid.UUID
	fetchedCompany *company.Company
	testUser       *user.User = user.NewUser(
		"company_test@gmail.com",
		"123456789",
		"+77555553535",
	)
)

func TestMain(m *testing.M) {
	logger.InitFileLogger("company_test.logs")
	code := m.Run()
	os.Exit(code)
}

func TestCompany(t *testing.T) {
	err := auth.Register(testUser.Login, testUser.PasswordHash, testUser.Phone)
	if err != nil {
		panic("failed prepare user for  test company")
	}
	fetchedUser, err := auth.Login(testUser.Login, testUser.PasswordHash)
	if err != nil {
		panic("failed prepare user for  test company")
	}
	userId = fetchedUser.ID

	t.Run("NewCompany", func(t *testing.T) {
		var err error
		companyId, err = comp.NewCompany(userId, name, description)
		if err != nil {
			t.Fatalf("Failed NewCompany: %v", err)
		}
	})

	t.Run("GetUserCompanies", func(t *testing.T) {
		fetchedCompanies, err := comp.GetUserCompanies(userId)
		if err != nil {
			t.Fatalf("Failed GetUserCompanies: %v", err)
		}

		if len(*fetchedCompanies) != 1 {
			t.Errorf("Expected 1, got %d\n", len(*fetchedCompanies))
		}
		for _, value := range *fetchedCompanies {
			if value.ID != companyId {
				t.Errorf("Expected company id %s, got %s", companyId, value.ID)
			}
			fetchedCompany = &value
		}
	})

	t.Run("UpdateCompany", func(t *testing.T) {
		fetchedCompany.Name = "updateMyCompany"
		fetchedCompany.Description = "updateMyCompany"
		err := comp.UpdateCompany(fetchedCompany, companyId, userId)
		if err != nil {
			t.Fatalf("Failed UpdateCompany: %v", err)
		}
	})

	t.Run("GetCompany", func(t *testing.T) {
		updatedCompany, err := comp.GetCompany(userId, companyId)
		if err != nil {
			t.Fatalf("Failed GetCompany: %v", err)
		}
		if updatedCompany.Name != "updateMyCompany" {
			t.Errorf("Expected updateMyCompany, got %s", updatedCompany.Name)
		}
	})

}
