package util

import "database/sql"

var LOCAL bool = false

func GetDatabase() (*sql.DB, error) {
	var db *sql.DB
	var err error

	if LOCAL {
		// LOCAL
		db, err = sql.Open("mysql", "isvacbanned:root@tcp(localhost:3306)/isvacbanned")
	} else {
		// PROD
		db, err = sql.Open("mysql", "b4efd0d0f3c600:a5e2c7d6@tcp(us-cdbr-east-02.cleardb.com:3306)/heroku_bace7cf727a523d")
	}

	return db, err
}
