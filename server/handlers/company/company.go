package company

import (
	"labyrinth/logic"
	"labyrinth/models/company"

	"github.com/google/uuid"
)

var bl *logic.BusinessLogic = logic.NewBusinessLogic()

type companyRegister struct {
	Name        string `json: "name"`
	Description string `json: "description"`
}

type companyData struct {
	Id   uuid.UUID `json: "id`
	Name string    `json: "name"`
}

func companyToCompanyResponse(comp *company.Company) companyData {
	return companyData{
		comp.ID,
		comp.Name,
	}
}
