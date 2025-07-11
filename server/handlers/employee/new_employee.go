package employee

import (
	"database/sql"
	"encoding/json"
	"errors"
	"labyrinth/logger"
	"labyrinth/models/employee"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (e EmployeeHandlers) NewEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. Проверка аутентификации пользователя
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok || userID == uuid.Nil {
		logger.NewErrMessage("Invalid user ID in context",
			zap.String("operation", "NewEmployeeHandler"),
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
			zap.String("operation", "NewEmployeeHandler"),
			zap.String("variable", "user_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// 3. Проверка соответствия user_id в пути и в контексте
	if userPathId != userID {
		logger.NewWarnMessage("User ID mismatch",
			zap.String("operation", "NewEmployeeHandler"),
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
			zap.String("operation", "NewEmployeeHandler"),
			zap.String("variable", "company_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid company ID format", http.StatusBadRequest)
		return
	}

	// 6. Парсинг тела запроса
	var requestData employee.Employee
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		logger.NewWarnMessage("Failed to decode request body",
			zap.String("operation", "NewEmployeeHandler"),
			zap.String("company_id", companyId.String()),
			zap.Error(err),
		)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 7. Валидация данных
	if requestData.UserID == uuid.Nil {
		logger.NewWarnMessage("Empty user ID in request",
			zap.String("operation", "NewEmployeeHandler"),
			zap.String("company_id", companyId.String()),
		)
		http.Error(w, "User ID cannot be empty", http.StatusBadRequest)
		return
	}

	if requestData.PositionID == uuid.Nil {
		logger.NewWarnMessage("Empty position ID in request",
			zap.String("operation", "NewEmployeeHandler"),
			zap.String("company_id", companyId.String()),
		)
		http.Error(w, "Position ID cannot be empty", http.StatusBadRequest)
		return
	}

	// 9. Создание сотрудника
	err = bl.Employee.NewEmployee(userID, requestData.UserID, companyId, requestData.PositionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.NewWarnMessage("User or company not found",
				zap.String("operation", "NewEmployeeHandler"),
				zap.String("target_user_id", requestData.UserID.String()),
				zap.String("company_id", companyId.String()),
			)
			http.Error(w, "User or company not found", http.StatusNotFound)
			return
		}

		logger.NewErrMessage("Failed to create employee",
			zap.String("operation", "NewEmployeeHandler"),
			zap.String("target_user_id", requestData.UserID.String()),
			zap.String("company_id", companyId.String()),
			zap.Error(err),
		)
		http.Error(w, "Failed to create employee", http.StatusInternalServerError)
		return
	}

	// 10. Формирование успешного ответа
	response := map[string]interface{}{
		"status":  "success",
		"message": "Employee created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.NewErrMessage("Failed to encode response",
			zap.String("operation", "NewEmployeeHandler"),
			zap.Error(err),
		)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	logger.NewInfoMessage("Employee created successfully",
		zap.String("operation", "NewEmployeeHandler"),
		zap.String("admin_user_id", userID.String()),
		zap.String("target_user_id", requestData.UserID.String()),
		zap.String("company_id", companyId.String()),
	)
}
