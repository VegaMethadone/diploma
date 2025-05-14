package depemployee

import (
	"encoding/json"
	"labyrinth/logger"
	"labyrinth/models/depemployee"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (d DepEmployeeHandlers) UpdateDepEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. Проверка аутентификации пользователя
	userID, ok := ctx.Value("id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		logger.NewErrMessage("Invalid user ID in context",
			zap.String("operation", "UpdateDepEmployeeHandler"),
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
			zap.String("operation", "UpdateDepEmployeeHandler"),
			zap.String("variable", "user_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// 3. Проверка соответствия user_id в пути и в контексте
	if userPathId != userID {
		logger.NewWarnMessage("User ID mismatch",
			zap.String("operation", "UpdateDepEmployeeHandler"),
			zap.String("context_user_id", userID.String()),
			zap.String("path_user_id", userPathId.String()),
		)
		http.Error(w, "Forbidden: user ID mismatch", http.StatusForbidden)
		return
	}

	// 4. Парсинг company_id и department_id из пути
	companyId, err := uuid.Parse(vars["company_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid company ID format",
			zap.String("operation", "UpdateDepEmployeeHandler"),
			zap.String("variable", "company_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid company ID format", http.StatusBadRequest)
		return
	}

	departmentId, err := uuid.Parse(vars["department_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid department ID format",
			zap.String("operation", "UpdateDepEmployeeHandler"),
			zap.String("variable", "department_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid department ID format", http.StatusBadRequest)
		return
	}

	// 6. Парсинг тела запроса
	var updatedDepemployee depemployee.DepartmentEmployee
	if err := json.NewDecoder(r.Body).Decode(&updatedDepemployee); err != nil {
		logger.NewWarnMessage("Failed to decode request body",
			zap.String("operation", "UpdateDepEmployeeHandler"),
			zap.String("department_id", departmentId.String()),
			zap.Error(err),
		)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 7. Валидация данных
	if updatedDepemployee.EmployeeID == uuid.Nil {
		logger.NewWarnMessage("Empty employee ID in request",
			zap.String("operation", "UpdateDepEmployeeHandler"),
			zap.String("department_id", departmentId.String()),
		)
		http.Error(w, "Employee ID cannot be empty", http.StatusBadRequest)
		return
	}

	if updatedDepemployee.PositionID == uuid.Nil {
		logger.NewWarnMessage("Empty position ID in request",
			zap.String("operation", "UpdateDepEmployeeHandler"),
			zap.String("department_id", departmentId.String()),
		)
		http.Error(w, "Position ID cannot be empty", http.StatusBadRequest)
		return
	}

	// 10. Обновление данных сотрудника в департаменте
	if err := bl.DepartmentEmployee.UpdateDepEmployee(updatedDepemployee.EmployeeID, departmentId, &updatedDepemployee); err != nil {
		logger.NewErrMessage("Failed to update department employee",
			zap.String("operation", "UpdateDepEmployeeHandler"),
			zap.String("employee_id", updatedDepemployee.EmployeeID.String()),
			zap.String("department_id", departmentId.String()),
			zap.Error(err),
		)
		http.Error(w, "Failed to update employee in department", http.StatusInternalServerError)
		return
	}

	// 11. Формирование успешного ответа
	response := map[string]interface{}{
		"status":  "success",
		"message": "Department employee updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.NewErrMessage("Failed to encode response",
			zap.String("operation", "UpdateDepEmployeeHandler"),
			zap.String("employee_id", updatedDepemployee.EmployeeID.String()),
			zap.Error(err),
		)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	logger.NewInfoMessage("Department employee updated successfully",
		zap.String("operation", "UpdateDepEmployeeHandler"),
		zap.String("admin_user_id", userID.String()),
		zap.String("company_id", companyId.String()),
		zap.String("department_id", departmentId.String()),
		zap.String("employee_id", updatedDepemployee.EmployeeID.String()),
		zap.String("new_position_id", updatedDepemployee.PositionID.String()),
	)
}
