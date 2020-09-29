package model

import (
	"isvacbanned/util"
)

type UserModel struct{}

// GetUserID returns database user id for a telegram user id
func (u UserModel) GetUserID(telegramID int) (int64, error) {
	db, err := util.GetDatabase()

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	row := db.QueryRow("SELECT id FROM user WHERE telegram_id = ?", telegramID)

	var userID int64

	err = row.Scan(&userID)

	return userID, err
}

// CreateUser inserts a new user in the database
func (u UserModel) CreateUser(firstName, username string, telegramID int) int64 {
	db, err := util.GetDatabase()

	stmt, err := db.Prepare("INSERT INTO user(first_name, username, telegram_id) VALUES(?, ?, ?)")

	if err != nil {
		panic(err.Error())
	}

	res, err := stmt.Exec(firstName, username, telegramID)

	defer stmt.Close()
	if err != nil {
		panic(err.Error())
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		panic(err.Error())
	}

	return lastID
}

//InactivateUser sets user flag is_active to false
func (u UserModel) InactivateUser(userID int64) int64 {
	db, err := util.GetDatabase()

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	stmt, err := db.Prepare("UPDATE user SET is_active = false WHERE id = ?")

	if err != nil {
		panic(err.Error())
	}

	res, err := stmt.Exec(userID)

	defer stmt.Close()
	if err != nil {
		panic(err.Error())
	}

	rows, err := res.LastInsertId()
	if err != nil {
		panic(err.Error())
	}

	return rows
}

//ActivateUser sets user flag is_active to true
func (u UserModel) ActivateUser(userID int64) int64 {
	db, err := util.GetDatabase()

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	stmt, err := db.Prepare("UPDATE user SET is_active = true WHERE id = ?")

	if err != nil {
		panic(err.Error())
	}

	res, err := stmt.Exec(userID)

	defer stmt.Close()
	if err != nil {
		panic(err.Error())
	}

	rows, err := res.LastInsertId()
	if err != nil {
		panic(err.Error())
	}

	return rows
}
