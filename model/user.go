package model

import "database/sql"

type UserModel struct{}

// GetUserID returns database user id for a telegram user id
func (u UserModel) GetUserID(telegramID int) (int64, error) {
	// PROD
	//db, err := sql.Open("mysql", "b4efd0d0f3c600:a5e2c7d6@tcp(us-cdbr-east-02.cleardb.com:3306)/heroku_bace7cf727a523d")
	// LOCAL
	db, err := sql.Open("mysql", "isvacbanned:root@tcp(localhost:3306)/isvacbanned")

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
	// PROD
	//db, err := sql.Open("mysql", "b4efd0d0f3c600:a5e2c7d6@tcp(us-cdbr-east-02.cleardb.com:3306)/heroku_bace7cf727a523d")
	// LOCAL
	db, err := sql.Open("mysql", "isvacbanned:root@tcp(localhost:3306)/isvacbanned")

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

func (u UserModel) InactivateUser(userID int64) int64 {
	// PROD
	//db, err := sql.Open("mysql", "b4efd0d0f3c600:a5e2c7d6@tcp(us-cdbr-east-02.cleardb.com:3306)/heroku_bace7cf727a523d")
	// LOCAL
	db, err := sql.Open("mysql", "isvacbanned:root@tcp(localhost:3306)/isvacbanned")

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

func (u UserModel) ActivateUser(userID int64) int64 {
	// PROD
	//db, err := sql.Open("mysql", "b4efd0d0f3c600:a5e2c7d6@tcp(us-cdbr-east-02.cleardb.com:3306)/heroku_bace7cf727a523d")
	// LOCAL
	db, err := sql.Open("mysql", "isvacbanned:root@tcp(localhost:3306)/isvacbanned")

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