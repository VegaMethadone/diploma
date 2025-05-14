package depposition

import (
	"encoding/json"
	"labyrinth/logger"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewDepPositionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. Проверка аутентификации пользователя
	userID, ok := ctx.Value("id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		logger.NewErrMessage("Invalid user ID in context",
			zap.String("operation", "NewDepPositionHandler"),
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
			zap.String("operation", "NewDepPositionHandler"),
			zap.String("variable", "user_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// 3. Проверка соответствия user_id в пути и в контексте
	if userPathId != userID {
		logger.NewWarnMessage("User ID mismatch",
			zap.String("operation", "NewDepPositionHandler"),
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
			zap.String("operation", "NewDepPositionHandler"),
			zap.String("variable", "company_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid company ID format", http.StatusBadRequest)
		return
	}

	departmentId, err := uuid.Parse(vars["department_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid department ID format",
			zap.String("operation", "NewDepPositionHandler"),
			zap.String("variable", "department_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid department ID format", http.StatusBadRequest)
		return
	}

	// 6. Парсинг тела запроса
	var requestData deppositionData
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		logger.NewWarnMessage("Failed to decode request body",
			zap.String("operation", "NewDepPositionHandler"),
			zap.String("department_id", departmentId.String()),
			zap.Error(err),
		)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 7. Валидация данных
	if strings.TrimSpace(requestData.Name) == "" {
		logger.NewWarnMessage("Empty position name",
			zap.String("operation", "NewDepPositionHandler"),
			zap.String("department_id", departmentId.String()),
		)
		http.Error(w, "Position name cannot be empty", http.StatusBadRequest)
		return
	}

	if requestData.Lvl < 0 {
		logger.NewWarnMessage("Invalid position level",
			zap.String("operation", "NewDepPositionHandler"),
			zap.String("department_id", departmentId.String()),
			zap.Int("level", requestData.Lvl),
		)
		http.Error(w, "Position level must be positive", http.StatusBadRequest)
		return
	}

	// 8. Создание новой позиции
	_, err = bl.DepartmentEmployeePosition.NewDepemployeePos(departmentId, requestData.Lvl, requestData.Name)
	if err != nil {
		logger.NewErrMessage("Failed to create department position",
			zap.String("operation", "NewDepPositionHandler"),
			zap.String("department_id", departmentId.String()),
			zap.String("position_name", requestData.Name),
			zap.Error(err),
		)
		http.Error(w, "Failed to create position", http.StatusInternalServerError)
		return
	}

	// 9. Формирование ответа
	response := map[string]interface{}{
		"status":        "success",
		"message":       "Department position created successfully",
		"name":          requestData.Name,
		"level":         requestData.Lvl,
		"department_id": departmentId.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.NewErrMessage("Failed to encode response",
			zap.String("operation", "NewDepPositionHandler"),
			zap.Error(err),
		)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	logger.NewInfoMessage("Department position created successfully",
		zap.String("operation", "NewDepPositionHandler"),
		zap.String("user_id", userID.String()),
		zap.String("department_id", departmentId.String()),
		zap.String("position_name", requestData.Name),
		zap.Int("position_level", requestData.Lvl),
	)
}
