package logic

import (
	"fmt"
	pscompany "labyrinth/database/postgres/pscompany"
	psemployee "labyrinth/database/postgres/psemployee"
	myJwt "labyrinth/jwt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func LoginCompany(userId, companyId uuid.UUID) (string, error) {
	/*
	   1. Проверить существование компании
	   2. Проверить существование юзера в компании
	   3. Выдать нужный jwt токен этого employee
	*/

	// 1. Проверка существования компании
	// companyExists, err := CheckCompany(companyId)
	companyExists, err := pscompany.CheckCompany(companyId)
	if err != nil {
		return "", fmt.Errorf("failed to check company: %w", err)
	}
	if !companyExists {
		return "", fmt.Errorf("company not found")
	}

	// 2. Получение сотрудника
	employee_, err := psemployee.GetEmployee(userId, companyId)
	if err != nil {
		return "", fmt.Errorf("failed to get employee: %w", err)
	}
	if employee_ == nil {
		return "", fmt.Errorf("employee not found in this company")
	}

	// 3. Генерация JWT токена
	settings := jwt.MapClaims{
		"id":         employee_.Id.String(),
		"userId":     employee_.UserId.String(),
		"companyId":  employee_.CompanyId.String(),
		"positionId": employee_.PositionId.String(),
	}

	token_ := myJwt.NewToken(settings)

	return token_, nil
}
