package employeelogic_test

import (
	"fmt"
	"labyrinth/logger"
	authlogic "labyrinth/logic/authLogic"
	companylogic "labyrinth/logic/companyLogic"
	employeelogic "labyrinth/logic/employeeLogic"
	"labyrinth/models/company"
	"labyrinth/models/employee"
	"labyrinth/models/user"
	"os"
	"testing"

	"github.com/google/uuid"
)

var (
	auth authlogic.Auth              = authlogic.NewAuth()
	comp companylogic.CompanyLogic   = companylogic.NewCompanyLogic()
	emp  employeelogic.EmployeeLogic = employeelogic.NewEmployeeLogic()
	// dep            departmentlogic.DepartmentLogic = departmentlogic.NewDepartmentLogic()
	fetchedCompany *company.Company
	fetchedUser    *user.User
	targetUser     *user.User
	testUser       *user.User = user.NewUser(
		"employee_test1@gmail.com",
		"123456789",
		"+77955553535",
	)
	newUser *user.User = user.NewUser(
		"employee_test2@gmail.com",
		"123456789111",
		"+77955563535",
	)
	positionId uuid.UUID
)

func setup() error {
	err := auth.Register(testUser.Login, testUser.PasswordHash, testUser.Phone)
	if err != nil {
		return err
	}
	err = auth.Register(newUser.Login, newUser.PasswordHash, newUser.Phone)
	if err != nil {
		return err
	}

	fetchedUser, err = auth.Login(testUser.Login, testUser.PasswordHash)
	if err != nil {
		return err
	}
	targetUser, err = auth.Login(newUser.Login, newUser.PasswordHash)
	if err != nil {
		return err
	}

	companyId, err := comp.NewCompany(fetchedUser.ID, "myEmployee", "myEmployee")
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
	logger.InitFileLogger("employee_test.logs")
	if err := setup(); err != nil {
		fmt.Printf("Test setup failed: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()
	os.Exit(code)
}

func TestEmployee(t *testing.T) {
	employeeId := uuid.Nil
	t.Run("GetEmployee", func(t *testing.T) {
		fetchedEmployee, err := emp.GetEmployee(fetchedUser.ID, fetchedCompany.ID)
		if err != nil {
			t.Fatalf("Failed GetEmployee: %v", err)
		}
		positionId = fetchedEmployee.PositionID
		employeeId = fetchedEmployee.ID
	})

	t.Run("", func(t *testing.T) {
		err := emp.NewEmployee(employeeId, targetUser.ID, fetchedCompany.ID, positionId)
		if err != nil {
			t.Fatalf("Failed NewEmployee: %v", err)
		}
	})

	var futureEmployee *employee.Employee
	t.Run("GetAllEmployee", func(t *testing.T) {
		result, err := emp.GetAllEmployee(fetchedCompany.ID)
		if err != nil {
			t.Fatalf("Failed  GetAllEmployee: %v", err)
		}

		if len(*result) != 2 {
			t.Errorf("Expected 2 employee, got %d", len(*result))
		}

		for _, value := range *result {
			if value.UserID == targetUser.ID {
				futureEmployee = &value
			}
		}
	})

	t.Run("UpdateEmployee", func(t *testing.T) {
		futureEmployee.IsOnline = false
		err := emp.UpdateEmployee(fetchedUser.ID, fetchedCompany.ID, futureEmployee)
		if err != nil {
			t.Fatalf("Failed  UpdateEmployee:  %v", err)
		}
	})
}
