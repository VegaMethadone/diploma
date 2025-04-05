package logic

import (
	"encoding/json"
	"fmt"
	pscompany "labyrinth/database/postgres/pscompany"

	"github.com/google/uuid"
)

func GetUserCompanies(userId uuid.UUID) (string, error) {
	userCompanies, err := pscompany.GetUserCompanies(userId)
	if err != nil {
		return "", err
	}

	jsonData, err := json.Marshal(userCompanies)
	if err != nil {
		return "", fmt.Errorf("failed to marshal user companies: %w", err)
	}

	return string(jsonData), nil
}
