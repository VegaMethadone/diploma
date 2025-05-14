package company

import (
	"encoding/json"
	"labyrinth/logger"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (c CompanyHandlers) GetAllCompaniesHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Получаем и проверяем userID из контекста
	ctx := r.Context()
	userID, ok := ctx.Value("id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		logger.NewErrMessage("Invalid user ID in context",
			zap.String("operation", "GetAllCompaniesHandler"),
			zap.Any("context_values", ctx.Value("id")),
		)
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// 2. Извлекаем и проверяем user_id из пути
	vars := mux.Vars(r)
	userPathId, err := uuid.Parse(vars["user_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid path variable",
			zap.String("operation", "GetAllCompaniesHandler"),
			zap.String("variable", "user_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// 3. Проверяем соответствие user_id из пути и контекста
	if userPathId != userID {
		logger.NewWarnMessage("User ID mismatch",
			zap.String("operation", "GetAllCompaniesHandler"),
			zap.String("context_user_id", userID.String()),
			zap.String("path_user_id", userPathId.String()),
		)
		http.Error(w, "Forbidden: user ID mismatch", http.StatusForbidden)
		return
	}

	// 4. Получаем компании пользователя
	fetchedCompanies, err := bl.Company.GetUserCompanies(userID)
	if err != nil {
		logger.NewErrMessage("Failed to get user companies",
			zap.String("operation", "GetAllCompaniesHandler"),
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 5. Преобразуем в формат ответа
	userCompanies := make([]companyData, 0, len(*fetchedCompanies))
	for _, company := range *fetchedCompanies {
		userCompanies = append(userCompanies, companyToCompanyResponse(&company))
	}

	// 6. Отправляем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":    "success",
		"companies": userCompanies,
		"count":     len(userCompanies),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.NewErrMessage("Failed to encode response",
			zap.String("operation", "GetAllCompaniesHandler"),
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
	}

	logger.NewInfoMessage("Successfully retrieved user companies",
		zap.String("operation", "GetAllCompaniesHandler"),
		zap.String("user_id", userID.String()),
		zap.Int("companies_count", len(userCompanies)),
	)
}
