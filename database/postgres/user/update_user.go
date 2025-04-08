package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/models/user"

	"github.com/lib/pq"
)

func UpdateUser(ctx context.Context, sharedTx *sql.Tx, u *user.User) error {
	if sharedTx == nil {
		return errors.New("start transaction before query")
	}
	query := `
        UPDATE users
        SET
            login = $1,
            email = $2,
            phone = $3,
            first_name = $4,
            last_name = $5,
            bio = $6,
            telegram_username = $7,
            avatar_url = $8,
            is_active = $9,
            is_staff = $10,
            updated_at = NOW()
        WHERE id = $11
    `

	result, err := sharedTx.ExecContext(
		ctx,
		query,
		u.Login,
		u.Email,
		u.Phone,
		u.FirstName,
		u.LastName,
		u.Bio,
		u.TelegramUsername,
		u.AvatarURL,
		u.IsActive,
		u.IsStaff,
		u.ID,
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
		return fmt.Errorf("failed to update user: %w", err)
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
