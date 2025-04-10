package position

import (
	"time"

	"github.com/google/uuid"
)

const (
	PositionLevelNone    = 0
	PositionLevelOwner   = 1
	PositionLevelAdmin   = 2
	PositionLevelDefault = 3
)

type Position struct {
	ID        uuid.UUID `json:"id"`
	CompanyID uuid.UUID `json:"company_id"`
	Lvl       int       `json:"lvl"` // 0-нет, 1-владелец, 2-админ, 3-работник
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
