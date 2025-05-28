package company

import (
	"database/sql"
	"encoding/json"
	"errors"
	"labyrinth/logger"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (c CompanyHandlers) UpdateCompanyProfileHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. Проверка аутентификации пользователя
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok || userID == uuid.Nil {
		logger.NewErrMessage("Invalid user ID in context",
			zap.String("operation", "UpdateCompanyProfileHandler"),
			zap.Any("context_values", ctx.Value(userIDKey)),
		)
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// 2. Парсинг user_id из пути
	vars := mux.Vars(r)
	userPathId, err := uuid.Parse(vars["user_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid path variable",
			zap.String("operation", "UpdateCompanyProfileHandler"),
			zap.String("variable", "user_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// 3. Проверка соответствия user_id в пути и в контексте
	if userPathId != userID {
		logger.NewWarnMessage("User ID mismatch",
			zap.String("operation", "UpdateCompanyProfileHandler"),
			zap.String("context_user_id", userID.String()),
			zap.String("path_user_id", userPathId.String()),
		)
		http.Error(w, "Forbidden: user ID mismatch", http.StatusForbidden)
		return
	}

	// 4. Парсинг company_id из пути
	companyId, err := uuid.Parse(vars["company_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid company ID format",
			zap.String("operation", "UpdateCompanyProfileHandler"),
			zap.String("variable", "company_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid company ID format", http.StatusBadRequest)
		return
	}

	// 5. Парсинг тела запроса
	var updatedCompany companyProfile
	if err := json.NewDecoder(r.Body).Decode(&updatedCompany); err != nil {
		logger.NewWarnMessage("Failed to decode request body",
			zap.String("operation", "UpdateCompanyProfileHandler"),
			zap.String("user_id", userID.String()),
			zap.String("company_id", companyId.String()),
			zap.Error(err),
		)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 6. Валидация данных
	if err := validateCompanyProfile(&updatedCompany); err != nil {
		logger.NewWarnMessage("Validation failed",
			zap.String("operation", "UpdateCompanyProfileHandler"),
			zap.String("user_id", userID.String()),
			zap.String("company_id", companyId.String()),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 7. Преобразование в доменную модель
	originCompany := cleanCompanyToCompany(companyId, userID, &updatedCompany)

	// 8. Обновление компании
	if err := bl.Company.UpdateCompany(originCompany, companyId, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.NewWarnMessage("Company not found",
				zap.String("operation", "UpdateCompanyProfileHandler"),
				zap.String("user_id", userID.String()),
				zap.String("company_id", companyId.String()),
			)
			http.Error(w, "Company not found", http.StatusNotFound)
			return
		}

		logger.NewErrMessage("Failed to update company",
			zap.String("operation", "UpdateCompanyProfileHandler"),
			zap.String("user_id", userID.String()),
			zap.String("company_id", companyId.String()),
			zap.Error(err),
		)
		http.Error(w, "Failed to update company", http.StatusInternalServerError)
		return
	}

	// 9. Формирование успешного ответа
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Company updated successfully",
	}); err != nil {
		logger.NewErrMessage("Failed to encode response",
			zap.String("operation", "UpdateCompanyProfileHandler"),
			zap.String("user_id", userID.String()),
			zap.String("company_id", companyId.String()),
			zap.Error(err),
		)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	logger.NewInfoMessage("Company profile updated successfully",
		zap.String("user_id", userID.String()),
		zap.String("company_id", companyId.String()),
	)
}

// Вспомогательная функция для валидации
func validateCompanyProfile(c *companyProfile) error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("company name cannot be empty")
	}
	if strings.TrimSpace(c.Description) == "" {
		return errors.New("company description cannot be empty")
	}
	return nil
}
