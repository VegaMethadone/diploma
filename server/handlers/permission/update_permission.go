package permission

import (
	"encoding/json"
	"labyrinth/logger"
	"labyrinth/notebook/models/permission"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (p PermissionHandlers) UpdatePermissionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok || userID == uuid.Nil {
		logger.NewErrMessage("Invalid user ID in context",
			zap.String("operation", "UpdatePermissionHandler"),
			zap.Any("context_values", ctx.Value(userIDKey)),
		)
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	userPathId, err := uuid.Parse(vars["user_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid path variable",
			zap.String("operation", "UpdatePermissionHandler"),
			zap.String("variable", "user_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	if userPathId != userID {
		logger.NewWarnMessage("User ID mismatch",
			zap.String("operation", "UpdatePermissionHandler"),
			zap.String("context_user_id", userID.String()),
			zap.String("path_user_id", userPathId.String()),
		)
		http.Error(w, "Forbidden: user ID mismatch", http.StatusForbidden)
		return
	}

	notebookId, err := uuid.Parse(vars["notebook_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid notebook ID",
			zap.String("operation", "UpdatePermissionHandler"),
			zap.String("variable", "notebook_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid notebook ID format", http.StatusBadRequest)
		return
	}

	var requestData permission.Permission
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		logger.NewWarnMessage("Failed to decode request body",
			zap.String("operation", "UpdatePermissionHandler"),
			zap.Error(err),
		)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = fsl.Permission.UpdatePermission(notebookId, &requestData)
	if err != nil {
		logger.NewErrMessage("Failed to update permission",
			zap.String("operation", "UpdatePermissionHandler"),
			zap.String("notebook_id", notebookId.String()),
			zap.Error(err),
		)

		http.Error(w, "Failed to update permission", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Notebook permission updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.NewErrMessage("Failed to encode response",
			zap.String("operation", "UpdatePermissionHandler"),
			zap.Error(err),
		)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
