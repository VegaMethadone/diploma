package logic

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

func GetUserCompanies(userId uuid.UUID) (string, error) {
	userCompanies, err := ps.GetUserCompanies(userId)
	if err != nil {
		return "", err
	}

	jsonData, err := json.Marshal(userCompanies)
	if err != nil {
		return "", fmt.Errorf("failed to marshal user companies: %w", err)
	}

	return string(jsonData), nil
}
