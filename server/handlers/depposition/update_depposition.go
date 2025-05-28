package depposition

import (
	"database/sql"
	"encoding/json"
	"errors"
	"labyrinth/logger"
	"labyrinth/models/depposition"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (d DepPositionHandlers) UpdateDepPositionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. Authentication check
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok || userID == uuid.Nil {
		logger.NewErrMessage("Invalid user ID in context",
			zap.String("operation", "UpdateDepPositionHandler"),
			zap.Any("context_values", ctx.Value(userIDKey)),
		)
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// 2. Parse path variables
	vars := mux.Vars(r)
	userPathId, err := uuid.Parse(vars["user_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid path variable",
			zap.String("operation", "UpdateDepPositionHandler"),
			zap.String("variable", "user_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// 3. Verify user ID match
	if userPathId != userID {
		logger.NewWarnMessage("User ID mismatch",
			zap.String("operation", "UpdateDepPositionHandler"),
			zap.String("context_user_id", userID.String()),
			zap.String("path_user_id", userPathId.String()),
		)
		http.Error(w, "Forbidden: user ID mismatch", http.StatusForbidden)
		return
	}

	// 4. Parse company and department IDs
	_, err = uuid.Parse(vars["company_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid company ID format",
			zap.String("operation", "UpdateDepPositionHandler"),
			zap.String("variable", "company_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid company ID format", http.StatusBadRequest)
		return
	}

	departmentId, err := uuid.Parse(vars["department_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid department ID format",
			zap.String("operation", "UpdateDepPositionHandler"),
			zap.String("variable", "department_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid department ID format", http.StatusBadRequest)
		return
	}

	positionId, err := uuid.Parse(vars["position_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid position ID format",
			zap.String("operation", "UpdateDepPositionHandler"),
			zap.String("variable", "position_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid position ID format", http.StatusBadRequest)
		return
	}

	// 6. Parse request body
	var updatedPos depposition.DepPosition
	if err := json.NewDecoder(r.Body).Decode(&updatedPos); err != nil {
		logger.NewWarnMessage("Failed to decode request body",
			zap.String("operation", "UpdateDepPositionHandler"),
			zap.String("position_id", positionId.String()),
			zap.Error(err),
		)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 7. Validate position data
	if strings.TrimSpace(updatedPos.Name) == "" {
		logger.NewWarnMessage("Empty position name",
			zap.String("operation", "UpdateDepPositionHandler"),
			zap.String("position_id", positionId.String()),
		)
		http.Error(w, "Position name cannot be empty", http.StatusBadRequest)
		return
	}

	if updatedPos.Level < 0 {
		logger.NewWarnMessage("Invalid position level",
			zap.String("operation", "UpdateDepPositionHandler"),
			zap.String("position_id", positionId.String()),
			zap.Int("level", updatedPos.Level),
		)
		http.Error(w, "Position level must be positive", http.StatusBadRequest)
		return
	}

	// 9. Update department position
	if err := bl.DepartmentEmployeePosition.UpdateDepEmployeePos(0, userID, departmentId, &updatedPos); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.NewWarnMessage("Position not found",
				zap.String("operation", "UpdateDepPositionHandler"),
				zap.String("position_id", positionId.String()),
			)
			http.Error(w, "Position not found", http.StatusNotFound)
			return
		}

		logger.NewErrMessage("Failed to update department position",
			zap.String("operation", "UpdateDepPositionHandler"),
			zap.String("position_id", positionId.String()),
			zap.Error(err),
		)
		http.Error(w, "Failed to update position", http.StatusInternalServerError)
		return
	}

	// 10. Prepare success response
	response := map[string]interface{}{
		"status":  "success",
		"message": "Department position updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.NewErrMessage("Failed to encode response",
			zap.String("operation", "UpdateDepPositionHandler"),
			zap.String("position_id", positionId.String()),
			zap.Error(err),
		)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	logger.NewInfoMessage("Department position updated successfully",
		zap.String("operation", "UpdateDepPositionHandler"),
		zap.String("user_id", userID.String()),
		zap.String("department_id", departmentId.String()),
		zap.String("position_id", positionId.String()),
		zap.String("new_name", updatedPos.Name),
		zap.Int("new_level", updatedPos.Level),
	)
}
