package company

import (
	"labyrinth/logic"
	"labyrinth/models/company"
	notebookLogic "labyrinth/notebook/logic"
	"time"

	"github.com/google/uuid"
)

var bl *logic.BusinessLogic = logic.NewBusinessLogic()
var fsl *notebookLogic.FileSystem = notebookLogic.NewFileSystem()

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

type companyProfile struct {
	Name        string    `json:"name"`         // Название компании
	Description string    `json:"description"`  // Описание
	LogoURL     string    `json:"logo_url"`     // Ссылка на логотип
	Industry    string    `json:"industry"`     // Отрасль
	IsVerified  bool      `json:"is_verified"`  // Подтверждена ли компания
	FoundedDate time.Time `json:"founded_date"` // Дата основания
	Address     string    `json:"address"`      // Адрес
	Phone       string    `json:"phone"`        // Телефон
	Email       string    `json:"email"`        // Email
}

func cleanCompanyToCompany(companyId, ownerId uuid.UUID, comp *companyProfile) *company.Company {
	return &company.Company{
		ID:          companyId,
		OwnerID:     ownerId,
		Name:        comp.Name,
		Description: comp.Description,
		LogoURL:     comp.LogoURL,
		Industry:    comp.Industry,
		Employees:   0,
		IsVerified:  comp.IsVerified,
		CreatedAt:   comp.FoundedDate,
		UpdatedAt:   time.Now(),
		FoundedDate: comp.FoundedDate,
		Address:     comp.Address,
		Phone:       comp.Phone,
		Email:       comp.Email,
		TaxNumber:   "",
	}
}

func companyToCleanCompany(comp *company.Company) *companyProfile {
	return &companyProfile{
		Name:        comp.Name,
		Description: comp.Description,
		LogoURL:     comp.LogoURL,
		Industry:    comp.Industry,
		IsVerified:  comp.IsVerified,
		FoundedDate: comp.CreatedAt,
		Address:     comp.Address,
		Phone:       comp.Phone,
		Email:       comp.Email,
	}
}
