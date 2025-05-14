package company

import (
	"database/sql"
	"encoding/json"
	"errors"
	"labyrinth/logger"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (c CompanyHandlers) GetCompanyHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. Проверка аутентификации пользователя
	userID, ok := ctx.Value("id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		logger.NewErrMessage("Invalid user ID in context",
			zap.String("operation", "GetCompanyHandler"),
			zap.Any("context_values", ctx.Value("id")),
		)
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// 2. Парсинг user_id из пути
	vars := mux.Vars(r)
	userPathId, err := uuid.Parse(vars["user_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid path variable",
			zap.String("operation", "GetCompanyHandler"),
			zap.String("variable", "user_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// 3. Проверка соответствия user_id в пути и в контексте
	if userPathId != userID {
		logger.NewWarnMessage("User ID mismatch",
			zap.String("operation", "GetCompanyHandler"),
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
			zap.String("operation", "GetCompanyHandler"),
			zap.String("variable", "company_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid company ID format", http.StatusBadRequest)
		return
	}

	// 5. Получение данных папки компании
	fetchedFolder, err := fsl.Folder.GetFolder(companyId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.NewWarnMessage("Company folder not found",
				zap.String("operation", "GetCompanyHandler"),
				zap.String("user_id", userID.String()),
				zap.String("company_id", companyId.String()),
			)
			http.Error(w, "Company folder not found", http.StatusNotFound)
			return
		}

		logger.NewErrMessage("Failed to get company folder",
			zap.String("operation", "GetCompanyHandler"),
			zap.String("user_id", userID.String()),
			zap.String("company_id", companyId.String()),
			zap.Error(err),
		)
		http.Error(w, "Failed to get company folder", http.StatusInternalServerError)
		return
	}

	// 6. Формирование успешного ответа
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "success",
		"folder":     fetchedFolder,
		"company_id": companyId.String(),
	}); err != nil {
		logger.NewErrMessage("Failed to encode response",
			zap.String("operation", "GetCompanyHandler"),
			zap.String("user_id", userID.String()),
			zap.String("company_id", companyId.String()),
			zap.Error(err),
		)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	logger.NewInfoMessage("Company folder retrieved successfully",
		zap.String("user_id", userID.String()),
		zap.String("company_id", companyId.String()),
	)
}
