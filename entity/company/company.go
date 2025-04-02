package company

import "github.com/google/uuid"

type CompanyLogin struct {
	UserId    uuid.UUID `json: "user_id"`
	CompanyId uuid.UUID `json: "company_id"`
}

type Company struct {
	Id          uuid.UUID `json: "id"`
	Owner       uuid.UUID `json: "owner"`
	Text        string    `json: "text"`
	Description string    `json: "description"`
}

func NewCompany(owner uuid.UUID, text, description string) *Company {
	return &Company{
		Id:          uuid.New(),
		Owner:       owner,
		Text:        text,
		Description: description,
	}
}
