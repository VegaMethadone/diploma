package department

import (
	"database/sql"
	"encoding/json"
	"errors"
	"labyrinth/logger"
	"labyrinth/models/department"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (d DepartmentHandlers) UpdateDepartmentProfileHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. Проверка аутентификации пользователя
	userID, ok := ctx.Value("id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		logger.NewErrMessage("Invalid user ID in context",
			zap.String("operation", "UpdateDepartmentProfileHandler"),
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
			zap.String("operation", "UpdateDepartmentProfileHandler"),
			zap.String("variable", "user_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// 3. Проверка соответствия user_id в пути и в контексте
	if userPathId != userID {
		logger.NewWarnMessage("User ID mismatch",
			zap.String("operation", "UpdateDepartmentProfileHandler"),
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
			zap.String("operation", "UpdateDepartmentProfileHandler"),
			zap.String("variable", "company_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid company ID format", http.StatusBadRequest)
		return
	}

	// 5. Парсинг department_id из пути
	departmentId, err := uuid.Parse(vars["department_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid department ID format",
			zap.String("operation", "UpdateDepartmentProfileHandler"),
			zap.String("variable", "department_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid department ID format", http.StatusBadRequest)
		return
	}

	// 7. Парсинг тела запроса
	var updatedDepartment department.Department
	if err := json.NewDecoder(r.Body).Decode(&updatedDepartment); err != nil {
		logger.NewWarnMessage("Failed to decode request body",
			zap.String("operation", "UpdateDepartmentProfileHandler"),
			zap.String("department_id", departmentId.String()),
			zap.Error(err),
		)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 8. Валидация данных
	if strings.TrimSpace(updatedDepartment.Name) == "" {
		logger.NewWarnMessage("Empty department name",
			zap.String("operation", "UpdateDepartmentProfileHandler"),
			zap.String("department_id", departmentId.String()),
		)
		http.Error(w, "Department name cannot be empty", http.StatusBadRequest)
		return
	}

	// 9. Установка ID департамента из пути
	updatedDepartment.ID = departmentId

	// 10. Обновление департамента
	if err := bl.Department.UpdateDepartment(userID, companyId, &updatedDepartment); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.NewWarnMessage("Department not found",
				zap.String("operation", "UpdateDepartmentProfileHandler"),
				zap.String("department_id", departmentId.String()),
			)
			http.Error(w, "Department not found", http.StatusNotFound)
			return
		}

		logger.NewErrMessage("Failed to update department",
			zap.String("operation", "UpdateDepartmentProfileHandler"),
			zap.String("user_id", userID.String()),
			zap.String("company_id", companyId.String()),
			zap.String("department_id", departmentId.String()),
			zap.Error(err),
		)
		http.Error(w, "Failed to update department", http.StatusInternalServerError)
		return
	}

	// 11. Формирование успешного ответа
	response := map[string]interface{}{
		"status":  "success",
		"message": "Department updated successfully",
		// "department_id": departmentId.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.NewErrMessage("Failed to encode response",
			zap.String("operation", "UpdateDepartmentProfileHandler"),
			zap.String("department_id", departmentId.String()),
			zap.Error(err),
		)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	logger.NewInfoMessage("Department updated successfully",
		zap.String("user_id", userID.String()),
		zap.String("company_id", companyId.String()),
		zap.String("department_id", departmentId.String()),
	)
}
