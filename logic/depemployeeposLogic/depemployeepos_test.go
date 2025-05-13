package depemployeeposlogic_test

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
	"labyrinth/models/depposition"
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
	fetchedDepartment *department.Department
	testUser          *user.User = user.NewUser(
		"departmentEmployeePos_test1@gmail.com",
		"1234567891000000",
		"+77855551111",
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

	fetched1User, err = auth.Login(testUser.Login, testUser.PasswordHash)
	if err != nil {
		return err
	}

	companyId, err := comp.NewCompany(fetched1User.ID, "myDepartmentPos", "myDepartmentPos")
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

	return nil
}

func TestMain(m *testing.M) {
	logger.InitFileLogger("department_employee_pos_test.logs")
	if err := setup(); err != nil {
		fmt.Printf("Test setup failed: %v\n", err)
		os.Exit(1)
	}
	code := m.Run()

	os.Exit(code)
}

func TestDepartmentPositition(t *testing.T) {
	positionId := uuid.Nil
	var err error
	t.Run("NewDepemployeePos", func(t *testing.T) {
		positionId, err = deppos.NewDepemployeePos(departmenId, 99, "nolifer")
		if err != nil {
			t.Fatalf("Failed NewDepemployeePos: %v", err)
		}
	})

	var newpos *depposition.DepPosition
	t.Run("GetAllDepEmployeePos", func(t *testing.T) {
		positions, err := deppos.GetAllDepEmployeePos(departmenId)
		if err != nil {
			t.Fatalf("Failed GetAllDepEmployeePos: %v", err)
		}

		for _, value := range *positions {
			if value.Id == positionId {
				newpos = &value
			}
		}
		if newpos == nil {
			t.Errorf("Expected position id: %s", positionId.String())
		}
	})

	t.Run("UpdateDepEmployeePos", func(t *testing.T) {
		employee, err := emp.GetEmployee(fetched1User.ID, fetchedCompany.ID)
		if err != nil {
			t.Fatalf("Failed UpdateDepEmployeePos => GetEmployee: %v", err)
		}
		newpos.Name = "NEWPOS_UPDATE"
		err = deppos.UpdateDepEmployeePos(0, employee.ID, departmenId, newpos)
		if err != nil {
			t.Fatalf("Failed UpdateDepEmployeePos: %v", err)
		}
	})
}
