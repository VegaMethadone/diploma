package user

import (
	"encoding/json"
	"labyrinth/logger"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (u UserHandlers) GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Получаем контекст и извлекаем userID
	ctx := r.Context()

	// userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok || userID == uuid.Nil {
		logger.NewErrMessage("Invalid user ID in context",
			zap.String("operation", "GetUserProfileHandler"),
			zap.Any("context_values", ctx.Value(userIDKey)),
		)
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// 3. Получаем профиль пользователя
	fetchedUser, err := bl.User.GetUserProfile(userID)
	if err != nil {
		logger.NewErrMessage("Failed to get user profile",
			zap.String("operation", "GetUserProfileHandler"),
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Подготавливаем ответ
	profile := newUserData(fetchedUser)

	// 5. Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(profile); err != nil {
		logger.NewErrMessage("Failed to encode response",
			zap.String("operation", "GetUserProfileHandler"),
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
		// Не отправляем повторную ошибку клиенту, так как заголовки уже записаны
	}

	logger.NewInfoMessage("User profile retrieved successfully",
		zap.String("operation", "GetUserProfileHandler"),
		zap.String("user_id", userID.String()),
	)
}
