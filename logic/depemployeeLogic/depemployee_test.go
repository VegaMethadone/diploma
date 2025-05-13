package depemployeelogic_test

import (
	"fmt"
	"labyrinth/logger"
	authlogic "labyrinth/logic/authLogic"
	companylogic "labyrinth/logic/companyLogic"
	departmentlogic "labyrinth/logic/departmentLogic"
	depemployeelogic "labyrinth/logic/depemployeeLogic"
	depemployeeposlogic "labyrinth/logic/depemployeeposLogic"
	employeelogic "labyrinth/logic/employeeLogic"
	"labyrinth/models/company"
	"labyrinth/models/department"
	"labyrinth/models/depemployee"
	"labyrinth/models/user"
	"os"
	"testing"

	"github.com/google/uuid"
)

var (
	auth              authlogic.Auth                         = authlogic.NewAuth()
	comp              companylogic.CompanyLogic              = companylogic.NewCompanyLogic()
	emp               employeelogic.EmployeeLogic            = employeelogic.NewEmployeeLogic()
	dep               departmentlogic.DepartmentLogic        = departmentlogic.NewDepartmentLogic()
	depemp            depemployeelogic.DepemployeeLogic      = depemployeelogic.NewDepemployeeLogic()
	deppos            depemployeeposlogic.DepemploeePosLogic = depemployeeposlogic.NewDepemploeePosLogic()
	fetchedCompany    *company.Company
	fetched1User      *user.User
	fetched2User      *user.User
	fetchedDepartment *department.Department
	testUser          *user.User = user.NewUser(
		"departmentEmployee_test1@gmail.com",
		"123456789",
		"+77855553665",
	)
	targetUser *user.User = user.NewUser(
		"departmentEmployee_test2@gmail.com",
		"12345678911",
		"+77855553661",
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
	err = auth.Register(targetUser.Login, targetUser.PasswordHash, targetUser.Phone)
	if err != nil {
		return err
	}

	fetched1User, err = auth.Login(testUser.Login, testUser.PasswordHash)
	if err != nil {
		return err
	}
	fetched2User, err = auth.Login(targetUser.Login, targetUser.PasswordHash)
	if err != nil {
		return err
	}

	companyId, err := comp.NewCompany(fetched1User.ID, "myDepartment", "myDepartment")
	if err != nil {
		return nil
	}

	fetchedCompany, err = comp.GetCompany(fetched1User.ID, companyId)
	if err != nil {
		return err
	}

	departmenId, depEmployeeId, depPositionId, err = dep.NewDepartment(fetched1User.ID, fetchedCompany.ID, fetchedCompany.ID, "myDepartmentEmployee", "myDepartmentEmployee")
	if err != nil {
		return err
	}

	fetched1Employee, err := emp.GetEmployee(fetched1User.ID, fetchedCompany.ID)
	if err != nil {
		return err
	}

	err = emp.NewEmployee(fetched1Employee.ID, fetched2User.ID, fetchedCompany.ID, fetched1Employee.PositionID)
	if err != nil {
		return err
	}

	return nil
}

func TestMain(m *testing.M) {
	logger.InitFileLogger("department_employee_test.logs")
	if err := setup(); err != nil {
		fmt.Printf("Test setup failed: %v\n", err)
		os.Exit(1)
	}
	code := m.Run()

	os.Exit(code)
}

func TestDepartment(t *testing.T) {
	employeeId := uuid.Nil
	t.Run("NewDepemployee", func(t *testing.T) {
		res, err := emp.GetEmployee(fetched2User.ID, fetchedCompany.ID)
		if err != nil {
			t.Fatalf("Failed NewDepemployee  => GetEmployee: %v", err)
		}
		employeeId = res.ID
		err = depemp.NewDepemployee(res.ID, departmenId, depPositionId)
		if err != nil {
			t.Fatalf("Failed NewDepemployee: %v", err)
		}
	})

	var fetchedDepEmployee *depemployee.DepartmentEmployee
	t.Run("GetAllDepEmployees", func(t *testing.T) {
		res, err := depemp.GetAllDepEmployees(departmenId)
		if err != nil {
			t.Fatalf("Failed GetAllDepEmployees: %v", err)
		}
		for _, value := range *res {
			if value.EmployeeID == employeeId {
				fetchedDepEmployee = &value
			}
		}
		if fetchedDepEmployee == nil {
			t.Errorf("Expected department employee with id: %s", employeeId.String())
		}
	})

	t.Run("UpdateDepEmployee", func(t *testing.T) {
		fetchedDepEmployee.IsActive = false
		err := depemp.UpdateDepEmployee(employeeId, departmenId, fetchedDepEmployee)
		if err != nil {
			t.Fatalf("Failed UpdateDepEmployee: %v", err)
		}
	})

	t.Run("GetDepartmentEmployee", func(t *testing.T) {
		res, err := depemp.GetDepartmentEmployee(employeeId, departmenId)
		if err != nil {
			t.Fatalf("Failed GetDepartmentEmployee: %v", err)
		}
		if res.IsActive {
			t.Errorf("Expected isActive false, gut true")
		}
	})
}
