package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"labyrinth/models/user"
)

func (p PostgresUser) GetUserByCredentials(
	ctx context.Context,
	sharedTx *sql.Tx,
	login,
	password string,
) (*user.User, error) {
	if sharedTx == nil {
		return nil, errors.New("start transaction before query")
	}

	query := `
        SELECT 
            id, login, password_hash, email, email_verified,
            phone, phone_verified, first_name, last_name,
            bio, telegram_username, avatar_url,
            created_at, updated_at, last_login_at,
            is_active, is_staff
        FROM users
        WHERE login = $1
        LIMIT 1
    `
	// query := `
	//     SELECT
	//         id, login, password_hash, email, email_verified,
	//         phone, phone_verified, first_name, last_name,
	//         bio, telegram_username, avatar_url,
	//         created_at, updated_at, last_login_at,
	//         is_active, is_staff
	//     FROM users
	//     WHERE login = $1 and password_hash = $2
	//     LIMIT 1
	// `

	var u user.User
	err := sharedTx.QueryRow(query, login).Scan(
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
		if err == sql.ErrNoRows {
			sharedTx.Rollback()
			return nil, fmt.Errorf("invalid credentials")
		}
		sharedTx.Rollback()
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &u, nil
}
