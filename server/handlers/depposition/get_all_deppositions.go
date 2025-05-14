package depposition

import (
	"encoding/json"
	"labyrinth/logger"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (d DepPositionHandlers) GetAllDepPositionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. Проверка аутентификации пользователя
	userID, ok := ctx.Value("id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		logger.NewErrMessage("Invalid user ID in context",
			zap.String("operation", "GetAllDepPositionHandler"),
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
			zap.String("operation", "GetAllDepPositionHandler"),
			zap.String("variable", "user_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// 3. Проверка соответствия user_id в пути и в контексте
	if userPathId != userID {
		logger.NewWarnMessage("User ID mismatch",
			zap.String("operation", "GetAllDepPositionHandler"),
			zap.String("context_user_id", userID.String()),
			zap.String("path_user_id", userPathId.String()),
		)
		http.Error(w, "Forbidden: user ID mismatch", http.StatusForbidden)
		return
	}

	// 4. Парсинг company_id и department_id из пути
	_, err = uuid.Parse(vars["company_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid company ID format",
			zap.String("operation", "GetAllDepPositionHandler"),
			zap.String("variable", "company_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid company ID format", http.StatusBadRequest)
		return
	}

	departmentId, err := uuid.Parse(vars["department_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid department ID format",
			zap.String("operation", "GetAllDepPositionHandler"),
			zap.String("variable", "department_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid department ID format", http.StatusBadRequest)
		return
	}

	// 6. Получение списка позиций департамента
	positions, err := bl.DepartmentEmployeePosition.GetAllDepEmployeePos(departmentId)
	if err != nil {
		logger.NewErrMessage("Failed to get department positions",
			zap.String("operation", "GetAllDepPositionHandler"),
			zap.String("department_id", departmentId.String()),
			zap.Error(err),
		)
		http.Error(w, "Failed to get department positions", http.StatusInternalServerError)
		return
	}

	// 7. Логирование успешного выполнения
	logger.NewInfoMessage("Department positions retrieved successfully",
		zap.String("operation", "GetAllDepPositionHandler"),
		zap.String("user_id", userID.String()),
		zap.String("department_id", departmentId.String()),
	)

	// 8. Формирование ответа
	response := map[string]interface{}{
		"status":    "success",
		"positions": positions,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.NewErrMessage("Failed to encode response",
			zap.String("operation", "GetAllDepPositionHandler"),
			zap.String("department_id", departmentId.String()),
			zap.Error(err),
		)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
