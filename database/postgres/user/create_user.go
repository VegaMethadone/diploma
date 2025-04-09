package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/models/user"

	"github.com/lib/pq"
)

func (p PostgresUser) CreateUser(ctx context.Context, sharedTx *sql.Tx, u *user.User) error {
	if sharedTx == nil {
		return errors.New("start transaction before query")
	}
	query := `
        INSERT INTO users (
            id,
            login,
            password_hash,
            email,
            phone,
            first_name,
            last_name,
            bio,
            telegram_username,
            avatar_url,
            email_verified,
            phone_verified,
            created_at,
            updated_at,
            last_login_at,
            is_active,
            is_staff
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
    `

	result, err := sharedTx.ExecContext(
		ctx,
		query,
		u.ID,
		u.Login,
		u.PasswordHash,
		u.Email,
		u.Phone,
		u.FirstName,
		u.LastName,
		u.Bio,
		u.TelegramUsername,
		u.AvatarURL,
		u.EmailVerified,
		u.PhoneVerified,
		u.CreatedAt,
		u.UpdatedAt,
		u.LastLoginAt,
		u.IsActive,
		u.IsStaff,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Constraint {
			case "users_login_key":
				return errors.New("login already exists")
			case "users_email_key":
				return errors.New("email already exists")
			}
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found (id: %s)", u.ID)
	}

	return nil
}
