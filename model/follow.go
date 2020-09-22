package model

import (
	"database/sql"
	"log"
	"strconv"

	//import mysql driver
	_ "github.com/go-sql-driver/mysql"
)

type UsersFollowed struct {
	ID      int
	SteamID string
}

// FollowSteamUser links a telegram user to a steam user which is being followed
func FollowSteamUser(chatID int64, steamID string) int64 {
	// PROD
	//db, err := sql.Open("mysql", "b4efd0d0f3c600:a5e2c7d6@tcp(us-cdbr-east-02.cleardb.com:3306)/heroku_bace7cf727a523d")
	// LOCAL
	db, err := sql.Open("mysql", "isvacbanned:root@tcp(localhost:3306)/isvacbanned")

	if err != nil {
		panic(err.Error())
	}

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

//GetAllIncompletedFollowedUsers get all fallowed steam user for every telkegram user
func GetAllIncompletedFollowedUsers() map[int64][]UsersFollowed {
	// PROD
	//db, err := sql.Open("mysql", "b4efd0d0f3c600:a5e2c7d6@tcp(us-cdbr-east-02.cleardb.com:3306)/heroku_bace7cf727a523d")
	// LOCAL
	db, err := sql.Open("mysql", "isvacbanned:root@tcp(localhost:3306)/isvacbanned")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	rows, err := db.Query("SELECT id, chat_id, steam_id FROM follow WHERE is_completed <> true")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	m := make(map[int64][]UsersFollowed)
	for rows.Next() {
		var (
			id      int
			chatID  int64
			steamID string
		)

		if err = rows.Scan(&id, &chatID, &steamID); err != nil {
			if err != sql.ErrNoRows {
				panic(err.Error())
			}

			return nil
		}

		m[chatID] = append(m[chatID], UsersFollowed{ID: id, SteamID: steamID})
	}

	return m

}

func SetFollowedUserToCompleted(id []int) int64 {
	// PROD
	//db, err := sql.Open("mysql", "b4efd0d0f3c600:a5e2c7d6@tcp(us-cdbr-east-02.cleardb.com:3306)/heroku_bace7cf727a523d")
	// LOCAL
	db, err := sql.Open("mysql", "isvacbanned:root@tcp(localhost:3306)/isvacbanned")

	if err != nil {
		panic(err.Error())
	}

	stmt, err := db.Prepare("UPDATE follow SET is_completed = true where id in(?)")

	if err != nil {
		panic(err.Error())
	}

	res, err := stmt.Exec(sliceToStringParam(id))

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

func sliceToStringParam(ids []int) string {
	if len(ids) == 0 {
		return ""
	}

	// Appr. 3 chars per num plus the comma.
	estimate := len(ids) * 4
	b := make([]byte, 0, estimate)
	// Or simply
	//   b := []byte{}
	for _, n := range ids {
		b = strconv.AppendInt(b, int64(n), 10)
		b = append(b, ',')
	}
	b = b[:len(b)-1]
	return string(b)
}
