package user

import (
	"encoding/json"
	"errors"
	"labyrinth/logger"
	"labyrinth/server/handlers/internal/halper"
	"net/http"
	"net/mail"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func UpdateUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Получаем контекст и извлекаем userID
	ctx := r.Context()

	userID, ok := ctx.Value("id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		logger.NewErrMessage("Invalid user ID in context",
			zap.String("operation", "UpdateUserProfileHandler"),
			zap.Any("context_values", ctx.Value("id")),
		)
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	userPathId, err := uuid.Parse(vars["user_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid path variable",
			zap.String("operation", "UpdateUserProfileHandler"),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if userPathId != userID {
		logger.NewWarnMessage("Wrong user id",
			zap.String("operation", "UpdateUserProfileHandler"),
		)
		http.Error(w, "Wrong user id", http.StatusBadRequest)
		return
	}

	// 2. Проверяем Content-Type
	if err := halper.CheckBodyContent(r); err != nil {
		logger.NewWarnMessage("Invalid content type",
			zap.String("operation", "UpdateUserProfileHandler"),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	// 3. Декодируем тело запроса
	var requestUser userData
	if err := json.NewDecoder(r.Body).Decode(&requestUser); err != nil {
		logger.NewWarnMessage("Failed to decode request body",
			zap.String("operation", "UpdateUserProfileHandler"),
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 4. Валидация данных
	if err := validateUserData(requestUser); err != nil {
		logger.NewWarnMessage("Validation failed",
			zap.String("operation", "UpdateUserProfileHandler"),
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 5. Подготавливаем обновленные данные
	updatedUser := NewUser(&requestUser, userID)

	// 6. Обновляем профиль
	err = bl.User.UpdateUserProfile(updatedUser)
	if err != nil {
		logger.NewErrMessage("Failed to update user profile",
			zap.String("operation", "UpdateUserProfileHandler"),
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 7. Успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"status":  "success",
		"message": "Profile updated successfully",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.NewErrMessage("Failed to encode response",
			zap.String("operation", "UpdateUserProfileHandler"),
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
	}

	logger.NewInfoMessage("User profile updated successfully",
		zap.String("operation", "UpdateUserProfileHandler"),
		zap.String("user_id", userID.String()),
	)
}

func validateUserData(data userData) error {
	if strings.TrimSpace(data.Email) == "" {
		return errors.New("email is required")
	}

	if _, err := mail.ParseAddress(data.Email); err != nil {
		return errors.New("invalid email format")
	}

	return nil
}
