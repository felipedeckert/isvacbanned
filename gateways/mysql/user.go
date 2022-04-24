package mysql

import (
	"context"
	"database/sql"
)

type User struct {
	*sql.DB
}

// GetUserID returns database user id for a telegram user id
func (u *User) GetUserID(ctx context.Context, telegramID int64) (int64, error) {
	query := `SELECT id FROM user WHERE telegram_id = ?`

	var userID int64
	err := u.QueryRowContext(ctx, query, telegramID).Scan(&userID)
	if err != nil {
		return -1, err
	}

	return userID, nil
}

// CreateUser inserts a new user in the database
func (u *User) CreateUser(ctx context.Context, firstName, username string, telegramID int64) (int64, error) {
	query := `INSERT INTO user(first_name, username, telegram_id) VALUES(?, ?, ?)`

	stmt, err := u.PrepareContext(ctx, query)
	if err != nil {
		return -1, err
	}
	res, err := stmt.ExecContext(ctx, firstName, username, telegramID)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	lastID, err := res.LastInsertId()

	if err != nil {
		return -1, err
	}

	return lastID, nil
}

//InactivateUser sets user flag is_active to false
func (u *User) InactivateUser(ctx context.Context, userID int64) (int64, error) {
	query := `UPDATE user SET is_active = false WHERE id = ?`

	stmt, err := u.PrepareContext(ctx, query)
	if err != nil {
		return -1, err
	}
	res, err := stmt.ExecContext(ctx, userID)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return -1, err
	}

	return rowsAffected, nil
}

//ActivateUser sets user flag is_active to true
func (u *User) ActivateUser(ctx context.Context, userID int64) (int64, error) {
	query := `UPDATE user SET is_active = true WHERE id = ?`

	stmt, err := u.PrepareContext(ctx, query)
	if err != nil {
		return -1, err
	}
	res, err := stmt.ExecContext(ctx, userID)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return -1, err
	}

	return rowsAffected, nil
}
