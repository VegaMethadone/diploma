package position

import (
	"encoding/json"
	"labyrinth/logger"
	"labyrinth/models/position"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (p PositionHandlers) UpdatePositionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. Проверка аутентификации пользователя
	userID, ok := ctx.Value("id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		logger.NewErrMessage("Invalid user ID in context",
			zap.String("operation", "UpdatePositionHandler"),
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
			zap.String("operation", "UpdatePositionHandler"),
			zap.String("variable", "user_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// 3. Проверка соответствия user_id в пути и в контексте
	if userPathId != userID {
		logger.NewWarnMessage("User ID mismatch",
			zap.String("operation", "UpdatePositionHandler"),
			zap.String("context_user_id", userID.String()),
			zap.String("path_user_id", userPathId.String()),
		)
		http.Error(w, "Forbidden: user ID mismatch", http.StatusForbidden)
		return
	}

	// 4. Парсинг company_id и position_id из пути
	companyId, err := uuid.Parse(vars["company_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid company ID format",
			zap.String("operation", "UpdatePositionHandler"),
			zap.String("variable", "company_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid company ID format", http.StatusBadRequest)
		return
	}

	positionId, err := uuid.Parse(vars["position_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid position ID format",
			zap.String("operation", "UpdatePositionHandler"),
			zap.String("variable", "position_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid position ID format", http.StatusBadRequest)
		return
	}

	// 5. Парсинг тела запроса
	var requestData position.Position
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		logger.NewWarnMessage("Failed to decode request body",
			zap.String("operation", "UpdatePositionHandler"),
			zap.String("position_id", positionId.String()),
			zap.Error(err),
		)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 6. Валидация данных
	if strings.TrimSpace(requestData.Name) == "" {
		logger.NewWarnMessage("Empty position name",
			zap.String("operation", "UpdatePositionHandler"),
			zap.String("position_id", positionId.String()),
		)
		http.Error(w, "Position name cannot be empty", http.StatusBadRequest)
		return
	}

	if requestData.Lvl < 0 {
		logger.NewWarnMessage("Invalid position level",
			zap.String("operation", "UpdatePositionHandler"),
			zap.String("position_id", positionId.String()),
			zap.Int("level", requestData.Lvl),
		)
		http.Error(w, "Position level must be positive", http.StatusBadRequest)
		return
	}

	// Устанавливаем ID из пути, чтобы избежать подмены
	requestData.ID = positionId

	// 7. Обновление позиции
	err = bl.Position.UpdatePosition(userID, companyId, &requestData)
	if err != nil {
		logger.NewErrMessage("Failed to update position",
			zap.String("operation", "UpdatePositionHandler"),
			zap.String("user_id", userID.String()),
			zap.String("company_id", companyId.String()),
			zap.String("position_id", positionId.String()),
			zap.Error(err),
		)
		http.Error(w, "Failed to update position", http.StatusInternalServerError)
		return
	}

	// 8. Формирование ответа
	response := map[string]interface{}{
		"status":  "success",
		"message": "Position updated successfully",
		// "position_id": positionId.String(),
		// "name":        requestData.Name,
		// "level":       requestData.Lvl,
		// "company_id":  companyId.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.NewErrMessage("Failed to encode response",
			zap.String("operation", "UpdatePositionHandler"),
			zap.Error(err),
		)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	logger.NewInfoMessage("Position updated successfully",
		zap.String("operation", "UpdatePositionHandler"),
		zap.String("user_id", userID.String()),
		zap.String("company_id", companyId.String()),
		zap.String("position_id", positionId.String()),
		zap.String("position_name", requestData.Name),
		zap.Int("position_level", requestData.Lvl),
	)
}
