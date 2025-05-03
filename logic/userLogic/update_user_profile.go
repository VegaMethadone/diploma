package userlogic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/database/postgres"
	"labyrinth/logger"
	"labyrinth/models/user"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (a Userlogic) UpdateUserProfile(userProfile *user.User) error {
	// 1. Валидация входных данных
	if userProfile.ID == uuid.Nil {
		logger.NewWarnMessage("user_id cannot be empty",
			zap.String("email", userProfile.Email),
			zap.Time("update_profile_time", time.Now()),
		)
		return errors.New("user_id cannot be empty")
	}

	if userProfile.PasswordHash == "" {
		logger.NewWarnMessage("user password cannot be empty",
			zap.String("email", userProfile.Email),
			zap.Time("update_profile_time", time.Now()),
		)
		return errors.New("user password cannot be empty")
	}

	// 2. Инициализация подключения к БД
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "UpdateUserProfile"),
			zap.String("email", userProfile.Email),
		)
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer db.Close()

	// 3. Настройка контекста с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 4. Начало транзакции (должна быть read-write для обновления)
	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false})
	if err != nil {
		logger.NewErrMessage("Transaction begin failed",
			zap.Error(err),
			zap.String("operation", "UpdateUserProfile"),
			zap.String("email", userProfile.Email),
		)
		return fmt.Errorf("transaction failed: %w", err)
	}

	// 5. Обновление профиля пользователя
	ps := postgres.NewPostgresDB()
	err = ps.User.UpdateUser(ctx, tx, userProfile)
	if err != nil {
		// Откат транзакции в случае ошибки
		tx.Rollback()

		if err == sql.ErrNoRows {
			logger.NewWarnMessage("No rows were updated - user not found",
				zap.String("email", userProfile.Email),
				zap.String("user_id", userProfile.ID.String()),
			)
			return errors.New("no such user")
		}

		logger.NewErrMessage("Failed to update user profile",
			zap.Error(err),
			zap.String("email", userProfile.Email),
			zap.String("user_id", userProfile.ID.String()),
		)
		return fmt.Errorf("failed to update user profile: %w", err)
	}

	// 7. Загрудаю  данные  в minio

	// 8. Коммит транзакции при успешном обновлении
	if err := tx.Commit(); err != nil {
		logger.NewErrMessage("Transaction commit failed",
			zap.Error(err),
			zap.String("operation", "UpdateUserProfile"),
			zap.String("email", userProfile.Email),
		)
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	// 9. Логирование успешного обновления
	logger.NewInfoMessage("User profile updated successfully",
		zap.String("user_id", userProfile.ID.String()),
		zap.String("email", userProfile.Email),
		zap.Time("update_time", time.Now()),
	)

	return nil
}
