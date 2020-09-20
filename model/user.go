package model

import "database/sql"

func GetUserID(telegramID int) (int64, error) {
	db, err := sql.Open("mysql", "b4efd0d0f3c600:a5e2c7d6@tcp(us-cdbr-east-02.cleardb.com:3306)/heroku_bace7cf727a523d")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	row := db.QueryRow("SELECT id FROM user WHERE telegram_id = ?", telegramID)

	var userID int64

	err = row.Scan(&userID)

	if err != nil {
		if err != sql.ErrNoRows {
			panic(err.Error())
		}

		return -1, err
	}

	return userID, nil
}

func CreateUser(firstName, username string, telegramID int64) int64 {
	db, err := sql.Open("mysql", "b4efd0d0f3c600:a5e2c7d6@tcp(us-cdbr-east-02.cleardb.com:3306)/heroku_bace7cf727a523d")

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
