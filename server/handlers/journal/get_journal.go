package journal

import (
	"encoding/json"
	"labyrinth/logger"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (j JournalHandler) GetNotebookHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. Проверка аутентификации пользователя
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok || userID == uuid.Nil {
		logger.NewErrMessage("Invalid user ID in context",
			zap.String("operation", "GetNotebookHandler"),
			zap.Any("context_values", ctx.Value(userIDKey)),
		)
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// 2. Парсинг параметров пути
	vars := mux.Vars(r)
	userPathId, err := uuid.Parse(vars["user_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid path variable",
			zap.String("operation", "GetNotebookHandler"),
			zap.String("variable", "user_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// 3. Проверка соответствия user_id в пути и в контексте
	if userPathId != userID {
		logger.NewWarnMessage("User ID mismatch",
			zap.String("operation", "GetNotebookHandler"),
			zap.String("context_user_id", userID.String()),
			zap.String("path_user_id", userPathId.String()),
		)
		http.Error(w, "Forbidden: user ID mismatch", http.StatusForbidden)
		return
	}

	notebookId, err := uuid.Parse(vars["notebook_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid notebook ID",
			zap.String("operation", "GetNotebookHandler"),
			zap.String("variable", "notebook_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid notebook ID format", http.StatusBadRequest)
		return
	}

	notebook, err := fsl.File.GetNotebook(notebookId)
	if err != nil {
		logger.NewErrMessage("Failed to get notebook",
			zap.String("operation", "GetNotebookHandler"),
			zap.String("notebook_id", notebookId.String()),
			zap.Error(err),
		)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status": "success",
		"data":   notebook,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.NewErrMessage("Failed to encode response",
			zap.String("operation", "GetNotebookHandler"),
			zap.Error(err),
		)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	logger.NewInfoMessage("Notebook retrieved successfully",
		zap.String("operation", "GetNotebookHandler"),
		zap.String("user_id", userID.String()),
		zap.String("notebook_id", notebookId.String()),
	)
}
