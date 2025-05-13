package departmentlogic_test

import (
	"fmt"
	"labyrinth/logger"
	authlogic "labyrinth/logic/authLogic"
	companylogic "labyrinth/logic/companyLogic"
	departmentlogic "labyrinth/logic/departmentLogic"
	"labyrinth/models/company"
	"labyrinth/models/department"
	"labyrinth/models/user"
	"os"
	"testing"

	"github.com/google/uuid"
)

var (
	auth              authlogic.Auth                  = authlogic.NewAuth()
	comp              companylogic.CompanyLogic       = companylogic.NewCompanyLogic()
	dep               departmentlogic.DepartmentLogic = departmentlogic.NewDepartmentLogic()
	fetchedCompany    *company.Company
	fetchedUser       *user.User
	fetchedDepartment *department.Department
	testUser          *user.User = user.NewUser(
		"department_test@gmail.com",
		"123456789",
		"+77855553535",
	)
	departmenId   uuid.UUID
	depEmployeeId uuid.UUID
	depPositionId uuid.UUID
)

func setup() error {
	err := auth.Register(testUser.Login, testUser.PasswordHash, testUser.Phone)
	if err != nil {
		return err
	}

	fetchedUser, err = auth.Login(testUser.Login, testUser.PasswordHash)
	if err != nil {
		return err
	}

	companyId, err := comp.NewCompany(fetchedUser.ID, "myDepartment", "myDepartment")
	if err != nil {
		return nil
	}

	fetchedCompany, err = comp.GetCompany(fetchedUser.ID, companyId)
	if err != nil {
		return err
	}

	return nil
}

func TestMain(m *testing.M) {
	logger.InitFileLogger("department_test.logs")
	if err := setup(); err != nil {
		fmt.Printf("Test setup failed: %v\n", err)
		os.Exit(1)
	}
	code := m.Run()

	os.Exit(code)
}

func TestDepartment(t *testing.T) {
	var err error
	t.Run("NewDepartment", func(t *testing.T) {
		departmenId, depEmployeeId, depPositionId, err = dep.NewDepartment(fetchedUser.ID, fetchedCompany.ID, fetchedCompany.ID, "myDepartment", "myDepartment")
		if err != nil {
			t.Fatalf("Failed NewDepartment: %v", err)
		}
	})
	t.Run("GetDepartment", func(t *testing.T) {
		fetchedDepartment, err = dep.GetDepartment(fetchedUser.ID, fetchedCompany.ID, departmenId)
		if err != nil {
			t.Fatalf("Failed GetDepartment: %v", err)
		}
	})

	t.Run("UpdateDepartment", func(t *testing.T) {
		fetchedDepartment.Name = "UpdatedmyDepartment"
		err = dep.UpdateDepartment(fetchedUser.ID, fetchedCompany.ID, fetchedDepartment)
		if err != nil {
			t.Fatalf("Failed UpdateDepartment: %v", err)
		}
	})

}
