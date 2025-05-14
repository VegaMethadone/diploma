package department

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

func (d DepartmentHandlers) NewDepartmentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. Проверка аутентификации пользователя
	userID, ok := ctx.Value("id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		logger.NewErrMessage("Invalid user ID in context",
			zap.String("operation", "NewDepartmentHandler"),
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
			zap.String("operation", "NewDepartmentHandler"),
			zap.String("variable", "user_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// 3. Проверка соответствия user_id в пути и в контексте
	if userPathId != userID {
		logger.NewWarnMessage("User ID mismatch",
			zap.String("operation", "NewDepartmentHandler"),
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
			zap.String("operation", "NewDepartmentHandler"),
			zap.String("variable", "company_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid company ID format", http.StatusBadRequest)
		return
	}

	// 5. Парсинг тела запроса
	var requestData departmentData
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		logger.NewWarnMessage("Failed to decode request body",
			zap.String("operation", "NewDepartmentHandler"),
			zap.String("user_id", userID.String()),
			zap.String("company_id", companyId.String()),
			zap.Error(err),
		)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 6. Валидация данных
	if strings.TrimSpace(requestData.Name) == "" {
		logger.NewWarnMessage("Empty department name",
			zap.String("operation", "NewDepartmentHandler"),
			zap.String("user_id", userID.String()),
			zap.String("company_id", companyId.String()),
		)
		http.Error(w, "Department name cannot be empty", http.StatusBadRequest)
		return
	}

	// 7. Получение данных сотрудника
	fetchedEmployee, err := bl.Employee.GetEmployee(userID, companyId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.NewWarnMessage("Employee not found",
				zap.String("operation", "NewDepartmentHandler"),
				zap.String("user_id", userID.String()),
				zap.String("company_id", companyId.String()),
			)
			http.Error(w, "Employee not found", http.StatusNotFound)
			return
		}

		logger.NewErrMessage("Failed to get employee",
			zap.String("operation", "NewDepartmentHandler"),
			zap.String("user_id", userID.String()),
			zap.String("company_id", companyId.String()),
			zap.Error(err),
		)
		http.Error(w, "Failed to get employee data", http.StatusInternalServerError)
		return
	}

	// 8. Создание департамента
	if _, _, _, err := bl.Department.NewDepartment(userID, companyId, requestData.ParentId, requestData.Name, requestData.Description); err != nil {
		logger.NewErrMessage("Failed to create department",
			zap.String("operation", "NewDepartmentHandler"),
			zap.String("user_id", userID.String()),
			zap.String("company_id", companyId.String()),
			zap.String("department_id", requestData.ParentId.String()),
			zap.Error(err),
		)
		http.Error(w, "Failed to create department", http.StatusInternalServerError)
		return
	}

	// 9. Создание папки департамента
	if err := fsl.Folder.CreateFolder(
		fetchedEmployee.ID,
		companyId,
		requestData.ParentId,
		requestData.ParentId,
		false,
		requestData.Name,
		requestData.Description,
	); err != nil {
		logger.NewErrMessage("Failed to create department folder",
			zap.String("operation", "NewDepartmentHandler"),
			zap.String("user_id", userID.String()),
			zap.String("company_id", companyId.String()),
			zap.String("department_id", requestData.ParentId.String()),
			zap.Error(err),
		)
		http.Error(w, "Failed to create department folder", http.StatusInternalServerError)
		return
	}

	// 10. Формирование успешного ответа
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"status":        "success",
		"message":       "Department created successfully",
		"department_id": requestData.ParentId.String(),
		"company_id":    companyId.String(),
	}); err != nil {
		logger.NewErrMessage("Failed to encode response",
			zap.String("operation", "NewDepartmentHandler"),
			zap.String("user_id", userID.String()),
			zap.String("company_id", companyId.String()),
			zap.String("department_id", requestData.ParentId.String()),
			zap.Error(err),
		)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	logger.NewInfoMessage("Department created successfully",
		zap.String("user_id", userID.String()),
		zap.String("company_id", companyId.String()),
		zap.String("department_id", requestData.ParentId.String()),
	)
}
