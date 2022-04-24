package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type User struct {
	*sql.DB
}

// GetUserID returns database user id for a telegram user id
func (u *User) GetUserID(ctx context.Context, telegramID int64) (int64, error) {
	query := `SELECT id FROM user WHERE telegram_id = $1`

	var userID int64
	err := u.QueryRowContext(ctx, query, telegramID).Scan(&userID)
	if err != nil {
		return -1, errors.New(fmt.Sprintf(`error get userID telegramID:%d`, telegramID))
	}

	return userID, nil
}

// CreateUser inserts a new user in the database
func (u *User) CreateUser(ctx context.Context, firstName, username string, telegramID int64) (int64, error) {
	query := `INSERT INTO user(first_name, username, telegram_id) VALUES($1, $2, $3)
				RETURNING  id`

	var lastID int64

	err := u.QueryRowContext(ctx, query, firstName, username, telegramID).Scan(&lastID)
	if err != nil {
		return -1, errors.New(fmt.Sprintf(`error creating user:%s`, firstName))
	}

	return lastID, nil
}

//InactivateUser sets user flag is_active to false
func (u *User) InactivateUser(ctx context.Context, userID int64) (int64, error) {
	query := `UPDATE user SET is_active = false WHERE id = $1
				RETURNING id`

	var id int64
	err := u.QueryRowContext(ctx, query, userID).Scan(&id)
	if err != nil {
		return -1, errors.New(fmt.Sprintf(`error inactivating user userID:%d`, userID))
	}

	return id, nil
}

//ActivateUser sets user flag is_active to true
func (u *User) ActivateUser(ctx context.Context, userID int64) (int64, error) {
	query := `UPDATE user SET is_active = true WHERE id = $1
				RETURNING id`

	var id int64
	err := u.QueryRowContext(ctx, query, userID).Scan(&id)
	if err != nil {
		return -1, errors.New(fmt.Sprintf(`error activating user userID:%d`, userID))
	}

	return id, nil
}
