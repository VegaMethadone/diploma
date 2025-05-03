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

func (u Userlogic) GetUserProfile(userId uuid.UUID) (*user.User, error) {
	if userId == uuid.Nil {
		logger.NewWarnMessage("",
			zap.String("operstion", "GetUserProfile"),
		)
		return nil, errors.New("userId is requred")
	}

	// 2. Инициализация подключения к БД
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "GetUserProfile"),
			zap.String("userId", userId.String()),
		)
		return nil, fmt.Errorf("database connection failed: %w", err)
	}
	defer db.Close()

	// 3. Настройка контекста с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 4. Начало транзакции (read-only)
	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		logger.NewErrMessage("Transaction begin failed",
			zap.Error(err),
			zap.String("operation", "login"),
			zap.String("userId", userId.String()),
		)
		return nil, fmt.Errorf("transaction failed: %w", err)
	}
	defer tx.Rollback() // Read-only транзакция - всегда Rollback

	// 5. Получаю данные пользователя
	ps := postgres.NewPostgresDB()
	fetchedUser, err := ps.User.GetUserByID(ctx, tx, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.NewWarnMessage("User not found",
				zap.String("userId", userId.String()),
			)
			return nil, errors.New("invalid user id")
		}

		logger.NewErrMessage("Failed to fetch user",
			zap.Error(err),
			zap.String("userId", userId.String()),
		)

		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	// 6. Вытаскиваю данные из minio если они были

	// 7. Аудит успешного входа
	logger.NewInfoMessage("User profile got successfully",
		zap.String("user_id", userId.String()),
		zap.Time("login_time", time.Now()),
	)

	return fetchedUser, nil
}
