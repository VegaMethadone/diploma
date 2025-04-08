package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/models/user"

	"github.com/google/uuid"
)

func GetUserByID(ctx context.Context, sharedTx *sql.Tx, id uuid.UUID) (*user.User, error) {
	if sharedTx == nil {
		return nil, errors.New("start transaction before query")
	}
	query := `
        SELECT 
            id,
            login,
            password_hash,
            email,
            email_verified,
            phone,
            phone_verified,
            first_name,
            last_name,
            bio,
            telegram_username,
            avatar_url,
            created_at,
            updated_at,
            last_login_at,
            is_active,
            is_staff
        FROM users
        WHERE id = $1
    `

	var u user.User
	err := sharedTx.QueryRowContext(ctx, query, id).Scan(
		&u.ID,
		&u.Login,
		&u.PasswordHash,
		&u.Email,
		&u.EmailVerified,
		&u.Phone,
		&u.PhoneVerified,
		&u.FirstName,
		&u.LastName,
		&u.Bio,
		&u.TelegramUsername,
		&u.AvatarURL,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.LastLoginAt,
		&u.IsActive,
		&u.IsStaff,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &u, nil
}
