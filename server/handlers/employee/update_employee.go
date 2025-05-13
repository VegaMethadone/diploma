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

func UpdateEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. Проверка аутентификации пользователя
	userID, ok := ctx.Value("id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		logger.NewErrMessage("Invalid user ID in context",
			zap.String("operation", "UpdateEmployeeHandler"),
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
			zap.String("operation", "UpdateEmployeeHandler"),
			zap.String("variable", "user_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// 3. Проверка соответствия user_id в пути и в контексте
	if userPathId != userID {
		logger.NewWarnMessage("User ID mismatch",
			zap.String("operation", "UpdateEmployeeHandler"),
			zap.String("context_user_id", userID.String()),
			zap.String("path_user_id", userPathId.String()),
		)
		http.Error(w, "Forbidden: user ID mismatch", http.StatusForbidden)
		return
	}

	// 4. Парсинг company_id и employee_id из пути
	companyId, err := uuid.Parse(vars["company_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid company ID format",
			zap.String("operation", "UpdateEmployeeHandler"),
			zap.String("variable", "company_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid company ID format", http.StatusBadRequest)
		return
	}

	employeeId, err := uuid.Parse(vars["employee_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid employee ID format",
			zap.String("operation", "UpdateEmployeeHandler"),
			zap.String("variable", "employee_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid employee ID format", http.StatusBadRequest)
		return
	}

	// 6. Парсинг тела запроса
	var requestData employee.Employee
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		logger.NewWarnMessage("Failed to decode request body",
			zap.String("operation", "UpdateEmployeeHandler"),
			zap.String("employee_id", employeeId.String()),
			zap.Error(err),
		)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 7. Валидация данных
	if requestData.PositionID == uuid.Nil {
		logger.NewWarnMessage("Empty position ID in request",
			zap.String("operation", "UpdateEmployeeHandler"),
			zap.String("employee_id", employeeId.String()),
		)
		http.Error(w, "Position ID cannot be empty", http.StatusBadRequest)
		return
	}

	// 10. Обновление данных сотрудника
	if err := bl.Employee.UpdateEmployee(userID, companyId, &requestData); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.NewWarnMessage("Employee not found",
				zap.String("operation", "UpdateEmployeeHandler"),
				zap.String("employee_id", employeeId.String()),
			)
			http.Error(w, "Employee not found", http.StatusNotFound)
			return
		}

		logger.NewErrMessage("Failed to update employee",
			zap.String("operation", "UpdateEmployeeHandler"),
			zap.String("employee_id", employeeId.String()),
			zap.Error(err),
		)
		http.Error(w, "Failed to update employee", http.StatusInternalServerError)
		return
	}

	// 11. Формирование успешного ответа
	response := map[string]interface{}{
		"status":      "success",
		"message":     "Employee updated successfully",
		"employee_id": employeeId.String(),
		"position_id": requestData.PositionID.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.NewErrMessage("Failed to encode response",
			zap.String("operation", "UpdateEmployeeHandler"),
			zap.String("employee_id", employeeId.String()),
			zap.Error(err),
		)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	logger.NewInfoMessage("Employee updated successfully",
		zap.String("operation", "UpdateEmployeeHandler"),
		zap.String("admin_user_id", userID.String()),
		zap.String("company_id", companyId.String()),
		zap.String("employee_id", employeeId.String()),
		zap.String("new_position_id", requestData.PositionID.String()),
	)
}
