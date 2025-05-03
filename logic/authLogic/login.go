package authlogic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/database/postgres"
	"labyrinth/logger"
	"labyrinth/models/user"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (a Auth) Login(mail, password string) (*user.User, error) {
	// 1. Валидация входных данных
	if mail == "" || password == "" {
		logger.NewWarnMessage("Empty credentials provided",
			zap.String("operation", "login"),
		)
		return nil, errors.New("email and password are required")
	}

	// 2. Инициализация подключения к БД
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Database connection failed",
			zap.Error(err),
			zap.String("operation", "login"),
			zap.String("email", mail),
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
			zap.String("email", mail),
		)
		return nil, fmt.Errorf("transaction failed: %w", err)
	}
	defer tx.Rollback() // Read-only транзакция - всегда Rollback

	// 5. Получение пользователя
	ps := postgres.NewPostgresDB()
	fetchedUser, err := ps.User.GetUserByCredentials(ctx, tx, mail, password)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.NewWarnMessage("User not found",
				zap.String("email", mail),
			)
			return nil, errors.New("invalid credentials")
		}

		logger.NewErrMessage("Failed to fetch user",
			zap.Error(err),
			zap.String("email", mail),
		)
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	// 6. Проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(fetchedUser.PasswordHash), []byte(password)); err != nil {
		logger.NewWarnMessage("Invalid password",
			zap.String("email", mail),
			zap.Error(err),
		)
		return nil, errors.New("invalid credentials")
	}

	// 7. Вытаскиваю данные из minio если они были

	// 8. Аудит успешного входа
	logger.NewInfoMessage("User logged in successfully",
		zap.String("user_id", fetchedUser.ID.String()),
		zap.String("email", mail),
		zap.Time("login_time", time.Now()),
	)

	return fetchedUser, nil
}
