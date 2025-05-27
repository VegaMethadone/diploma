package journal

import (
	"encoding/json"
	"labyrinth/logger"
	"labyrinth/notebook/models/journal"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (j JournalHandler) UpdateNotebookHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value("id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		logger.NewErrMessage("Invalid user ID in context",
			zap.String("operation", "UpdateNotebookHandler"),
			zap.Any("context_values", ctx.Value("id")),
		)
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	userPathId, err := uuid.Parse(vars["user_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid path variable",
			zap.String("operation", "UpdateNotebookHandler"),
			zap.String("variable", "user_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	if userPathId != userID {
		logger.NewWarnMessage("User ID mismatch",
			zap.String("operation", "UpdateNotebookHandler"),
			zap.String("context_user_id", userID.String()),
			zap.String("path_user_id", userPathId.String()),
		)
		http.Error(w, "Forbidden: user ID mismatch", http.StatusForbidden)
		return
	}

	notebookId, err := uuid.Parse(vars["notebook_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid notebook ID",
			zap.String("operation", "UpdateNotebookHandler"),
			zap.String("variable", "notebook_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid notebook ID format", http.StatusBadRequest)
		return
	}

	var requestData journal.Notebook
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		logger.NewWarnMessage("Failed to decode request body",
			zap.String("operation", "UpdateNotebookHandler"),
			zap.String("notebook_id", notebookId.String()),
			zap.Error(err),
		)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if strings.TrimSpace(requestData.UuidID) == "" {
		logger.NewWarnMessage("Empty notebook uuid in update request",
			zap.String("operation", "UpdateNotebookHandler"),
		)
		http.Error(w, "Notebook uuid cannot be empty", http.StatusBadRequest)
		return
	}

	if err := fsl.File.UpdateNotebook(notebookId, &requestData); err != nil {
		logger.NewErrMessage("Failed to update notebook",
			zap.String("operation", "UpdateNotebookHandler"),
			zap.String("notebook_id", notebookId.String()),
			zap.Error(err),
		)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Notebook updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.NewErrMessage("Failed to encode response",
			zap.String("operation", "UpdateNotebookHandler"),
			zap.Error(err),
		)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	logger.NewInfoMessage("Notebook updated successfully",
		zap.String("operation", "UpdateNotebookHandler"),
		zap.String("user_id", userID.String()),
		zap.String("notebook_id", notebookId.String()),
	)
}
