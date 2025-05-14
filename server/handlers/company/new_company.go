package company

import (
	"encoding/json"
	"errors"
	"labyrinth/logger"
	"labyrinth/server/handlers/internal/halper"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (c CompanyHandlers) NewCompanyHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Получаем и проверяем userID из контекста
	ctx := r.Context()
	userID, ok := ctx.Value("id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		logger.NewErrMessage("Invalid user ID in context",
			zap.String("operation", "NewCompanyHandler"),
			zap.Any("context_values", ctx.Value("id")),
		)
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// 2. Извлекаем и проверяем user_id из пути
	vars := mux.Vars(r)
	userPathId, err := uuid.Parse(vars["user_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid path variable",
			zap.String("operation", "NewCompanyHandler"),
			zap.String("variable", "user_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// 3. Проверяем соответствие user_id из пути и контекста
	if userPathId != userID {
		logger.NewWarnMessage("User ID mismatch",
			zap.String("operation", "NewCompanyHandler"),
			zap.String("context_user_id", userID.String()),
			zap.String("path_user_id", userPathId.String()),
		)
		http.Error(w, "Forbidden: user ID mismatch", http.StatusForbidden)
		return
	}

	// 4. Проверяем Content-Type
	if err := halper.CheckBodyContent(r); err != nil {
		logger.NewWarnMessage("Invalid content type",
			zap.String("operation", "NewCompanyHandler"),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	// 5. Декодируем тело запроса
	var requestData companyRegister
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		logger.NewWarnMessage("Failed to decode request body",
			zap.String("operation", "NewCompanyHandler"),
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 6. Валидация данных компании
	if err := validateCompanyData(requestData); err != nil {
		logger.NewWarnMessage("Company data validation failed",
			zap.String("operation", "NewCompanyHandler"),
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 7. Создаем компанию через бизнес-логику
	_, err = bl.Company.NewCompany(userID, requestData.Name, requestData.Description)
	if err != nil {
		logger.NewErrMessage("Failed to create company",
			zap.String("operation", "NewCompanyHandler"),
			zap.String("user_id", userID.String()),
			zap.String("company_name", requestData.Name),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 8. Успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"status":  "success",
		"message": "Company created successfully",
		"company": map[string]string{
			"name":        requestData.Name,
			"description": requestData.Description,
		},
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.NewErrMessage("Failed to encode response",
			zap.String("operation", "NewCompanyHandler"),
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
	}

	logger.NewInfoMessage("Company created successfully",
		zap.String("operation", "NewCompanyHandler"),
		zap.String("user_id", userID.String()),
		zap.String("company_name", requestData.Name),
	)
}

// Валидация данных компании
func validateCompanyData(data companyRegister) error {
	if strings.TrimSpace(data.Name) == "" {
		return errors.New("company name is required")
	}

	if len(data.Name) > 100 {
		return errors.New("company name is too long (max 100 chars)")
	}

	if len(data.Description) > 500 {
		return errors.New("company description is too long (max 500 chars)")
	}

	return nil
}
