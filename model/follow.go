package model

import (
	"database/sql"
	"log"
)

func FollowSteamUser(chatID int64, steamID string) int64 {
	db, err := sql.Open("mysql", "b4efd0d0f3c600:a5e2c7d6@tcp(us-cdbr-east-02.cleardb.com:3306)/heroku_bace7cf727a523d")

	stmt, err := db.Prepare("INSERT INTO follow(chat_id, steam_id) VALUES(?, ?)")

	if err != nil {
		panic(err.Error())
	}

	res, err := stmt.Exec(chatID, steamID)

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

func GetAllFollownUsers() map[int64][]string {
	db, err := sql.Open("mysql", "b4efd0d0f3c600:a5e2c7d6@tcp(us-cdbr-east-02.cleardb.com:3306)/heroku_bace7cf727a523d")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	rows, err := db.Query("SELECT chat_id, steam_id FROM follow")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	m := make(map[int64][]string)
	for rows.Next() {
		var (
			chatID  int64
			steamID string
		)

		if err = rows.Scan(&chatID, &steamID); err != nil {
			if err != sql.ErrNoRows {
				panic(err.Error())
			}

			return nil
		}

		m[chatID] = append(m[chatID], steamID)
	}

	return m

}
