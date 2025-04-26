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

func (a Auth) Register(mail, hashPassword, phone string) error {
	// 0. Валидация входных данных
	if mail == "" || hashPassword == "" || phone == "" {
		logger.NewWarnMessage("Empty credentials provided",
			zap.String("operation", "login"),
		)
		return errors.New("email and password and phone are required")
	}

	// 1. Инициализация подключения к БД
	db, err := sql.Open("postgres", postgres.GetConnection())
	if err != nil {
		logger.NewErrMessage("Failed to connect to database",
			zap.Error(err),
			zap.String("operation", "register"),
		)
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer db.Close()

	// 2. Настройка контекста с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 3. Начало транзакции
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		logger.NewErrMessage("Failed to begin transaction",
			zap.Error(err),
			zap.String("operation", "register"),
		)
		return fmt.Errorf("transaction failed: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 4. Проверка телефона
	ps := postgres.NewPostgresDB()
	exists, err := ps.User.CheckPhone(ctx, tx, phone)
	if err != nil {
		logger.NewErrMessage("Phone check failed",
			zap.Error(err),
			zap.String("phone", phone),
		)
		return fmt.Errorf("phone check failed: %w", err)
	}
	if exists {
		logger.NewWarnMessage("Phone already in use",
			zap.String("phone", phone),
		)
		return errors.New("phone number is already registered")
	}

	// 5. Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(hashPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.NewErrMessage("Password hashing failed",
			zap.Error(err),
		)
		return fmt.Errorf("password hashing failed: %w", err)
	}

	// 6. Создание пользователя
	newUser := user.NewUser(mail, string(hashedPassword), phone)
	newUser.ID, err = ps.UuidValidation.CheckAndReserveUUID(ctx, tx)
	if err != nil {
		logger.NewErrMessage("UUID generation failed",
			zap.Error(err),
		)
		return fmt.Errorf("UUID generation failed: %w", err)
	}

	if err := ps.User.CreateUser(ctx, tx, newUser); err != nil {
		logger.NewErrMessage("User creation failed",
			zap.Error(err),
			zap.String("email", mail),
			zap.String("phone", phone),
		)
		return fmt.Errorf("user creation failed: %w", err)
	}

	// 7. Фиксация транзакции
	if err := tx.Commit(); err != nil {
		logger.NewErrMessage("Transaction commit failed",
			zap.Error(err),
		)
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	// 8. Логирование успешной регистрации
	logger.NewInfoMessage("New user registered",
		zap.String("email", mail),
		zap.String("phone", phone),
		zap.Time("registered_at", time.Now()),
	)

	return nil
}
