package journal

import (
	"encoding/json"
	"labyrinth/logger"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (j JournalHandler) NewNotebookHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. Проверка аутентификации пользователя
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok || userID == uuid.Nil {
		logger.NewErrMessage("Invalid user ID in context",
			zap.String("operation", "NewNotebookHandler"),
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
			zap.String("operation", "NewNotebookHandler"),
			zap.String("variable", "user_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// 3. Проверка соответствия user_id в пути и в контексте
	if userPathId != userID {
		logger.NewWarnMessage("User ID mismatch",
			zap.String("operation", "NewNotebookHandler"),
			zap.String("context_user_id", userID.String()),
			zap.String("path_user_id", userPathId.String()),
		)
		http.Error(w, "Forbidden: user ID mismatch", http.StatusForbidden)
		return
	}

	companyId, err := uuid.Parse(vars["company_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid company ID",
			zap.String("operation", "NewNotebookHandler"),
			zap.String("variable", "company_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid company ID format", http.StatusBadRequest)
		return
	}

	departmentId, err := uuid.Parse(vars["department_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid department ID",
			zap.String("operation", "NewNotebookHandler"),
			zap.String("variable", "department_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid department ID format", http.StatusBadRequest)
		return
	}

	// 4. Парсинг тела запроса
	var requestData notebookRequest
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		logger.NewWarnMessage("Failed to decode request body",
			zap.String("operation", "NewNotebookHandler"),
			zap.String("company_id", companyId.String()),
			zap.Error(err),
		)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 5. Валидация данных
	if strings.TrimSpace(requestData.Name) == "" {
		logger.NewWarnMessage("Empty notebook name",
			zap.String("operation", "NewNotebookHandler"),
		)
		http.Error(w, "Notebook name cannot be empty", http.StatusBadRequest)
		return
	}

	// 6. Создание блокнота
	if err := fsl.File.NewNotebook(userID, companyId, departmentId, requestData.Name, requestData.Description); err != nil {
		logger.NewErrMessage("Failed to create notebook",
			zap.String("operation", "NewNotebookHandler"),
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 7. Формирование ответа
	response := map[string]interface{}{
		"status":  "success",
		"message": "Notebook created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.NewErrMessage("Failed to encode response",
			zap.String("operation", "NewNotebookHandler"),
			zap.Error(err),
		)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	logger.NewInfoMessage("Notebook created successfully",
		zap.String("operation", "NewNotebookHandler"),
		zap.String("user_id", userID.String()),
		zap.String("company_id", companyId.String()),
		zap.String("department_id", departmentId.String()),
		zap.String("notebook_name", requestData.Name),
	)
}
