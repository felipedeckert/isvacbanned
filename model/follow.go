package model

import (
	"database/sql"
	"log"
	"strconv"

	//import mysql driver
	_ "github.com/go-sql-driver/mysql"
)

type FollowModel struct{}

type UsersFollowed struct {
	ID           int
	SteamID      string
	OldNickname  string
	CurrNickname string
	IsCompleted  bool
}

// FollowSteamUser links a telegram user to a steam user which is being followed
func (f *FollowModel) FollowSteamUser(chatID int64, steamID, currNickname string, userID int64) int64 {
	// PROD
	//db, err := sql.Open("mysql", "b4efd0d0f3c600:a5e2c7d6@tcp(us-cdbr-east-02.cleardb.com:3306)/heroku_bace7cf727a523d")
	// LOCAL
	db, err := sql.Open("mysql", "isvacbanned:root@tcp(localhost:3306)/isvacbanned")

	if err != nil {
		panic(err.Error())
	}

	stmt, err := db.Prepare("INSERT INTO follow(chat_id, steam_id, user_id, old_nickname, curr_nickname) VALUES(?, ?, ?, ?, ?)")

	if err != nil {
		panic(err.Error())
	}

	res, err := stmt.Exec(chatID, steamID, userID, currNickname, currNickname)

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

// FollowSteamUser links a telegram user to a steam user which is being followed
func (f *FollowModel) UnfollowSteamUser(steamID string) int64 {
	// PROD
	//db, err := sql.Open("mysql", "b4efd0d0f3c600:a5e2c7d6@tcp(us-cdbr-east-02.cleardb.com:3306)/heroku_bace7cf727a523d")
	// LOCAL
	db, err := sql.Open("mysql", "isvacbanned:root@tcp(localhost:3306)/isvacbanned")

	if err != nil {
		panic(err.Error())
	}

	stmt, err := db.Prepare("UPDATE follow SET is_active = false where steam_id = ?")

	if err != nil {
		panic(err.Error())
	}

	res, err := stmt.Exec(steamID)

	defer stmt.Close()
	if err != nil {
		panic(err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		panic(err.Error())
	}

	return rows
}

//GetFollowerCountBySteamID get the number of followers of a steam player
func (f *FollowModel) GetFollowerCountBySteamID(steamID string) (int64, error) {
	// PROD
	//db, err := sql.Open("mysql", "b4efd0d0f3c600:a5e2c7d6@tcp(us-cdbr-east-02.cleardb.com:3306)/heroku_bace7cf727a523d")
	// LOCAL
	db, err := sql.Open("mysql", "isvacbanned:root@tcp(localhost:3306)/isvacbanned")

	if err != nil {
		panic(err.Error())
	}

	row := db.QueryRow(
		"SELECT COUNT(f.id) as count"+
			"	FROM follow f "+
			"	WHERE f.steam_id = ?", steamID)

	var count int64

	err = row.Scan(&count)

	return count, err
}

//GetAllIncompletedFollowedUsers get all fallowed steam user for every telegram user
func (f *FollowModel) GetAllIncompletedFollowedUsers() map[int64][]UsersFollowed {
	// PROD
	//db, err := sql.Open("mysql", "b4efd0d0f3c600:a5e2c7d6@tcp(us-cdbr-east-02.cleardb.com:3306)/heroku_bace7cf727a523d")
	// LOCAL
	db, err := sql.Open("mysql", "isvacbanned:root@tcp(localhost:3306)/isvacbanned")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	rows, err := db.Query("SELECT id, chat_id, steam_id, old_nickname, curr_nickname FROM follow WHERE is_completed <> true AND is_active = true")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	m := make(map[int64][]UsersFollowed)
	for rows.Next() {
		var (
			id           int
			chatID       int64
			steamID      string
			oldNickname  string
			currNickName string
		)

		if err = rows.Scan(&id, &chatID, &steamID, &oldNickname, &currNickName); err != nil {
			if err != sql.ErrNoRows {
				panic(err.Error())
			}

			return nil
		}

		m[chatID] = append(m[chatID], UsersFollowed{ID: id, SteamID: steamID, OldNickname: oldNickname, CurrNickname: currNickName})
	}

	return m

}

//GetUsersFollowed gets a slice os the nicknames (the old ones) of players followed by a user
func (f *FollowModel) GetUsersFollowed(userID int64) []UsersFollowed {
	// PROD
	//db, err := sql.Open("mysql", "b4efd0d0f3c600:a5e2c7d6@tcp(us-cdbr-east-02.cleardb.com:3306)/heroku_bace7cf727a523d")
	// LOCAL
	db, err := sql.Open("mysql", "isvacbanned:root@tcp(localhost:3306)/isvacbanned")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	rows, err := db.Query("SELECT old_nickname, is_completed FROM follow WHERE user_id = ?", userID)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	s := make([]UsersFollowed, 0)
	for rows.Next() {
		var (
			oldNickname string
			isCompleted bool
		)

		if err = rows.Scan(&oldNickname, &isCompleted); err != nil {
			if err != sql.ErrNoRows {
				panic(err.Error())
			}

			return nil
		}

		s = append(s, UsersFollowed{OldNickname: oldNickname, IsCompleted: isCompleted})
	}

	return s
}

func (f *FollowModel) SetCurrNickname(userId int, sanitizedActualNickname string) {
	// PROD
	//db, err := sql.Open("mysql", "b4efd0d0f3c600:a5e2c7d6@tcp(us-cdbr-east-02.cleardb.com:3306)/heroku_bace7cf727a523d")
	// LOCAL
	db, err := sql.Open("mysql", "isvacbanned:root@tcp(localhost:3306)/isvacbanned")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	stmt, err := db.Prepare("UPDATE follow SET curr_nickname = ? where id = ?")

	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(sanitizedActualNickname, userId)

	defer stmt.Close()
	if err != nil {
		panic(err.Error())
	}

}

// SetFollowedUserToCompleted sets a player status to completed, and it will not be followed anymore
func (f *FollowModel) SetFollowedUserToCompleted(id []int) int64 {
	// PROD
	//db, err := sql.Open("mysql", "b4efd0d0f3c600:a5e2c7d6@tcp(us-cdbr-east-02.cleardb.com:3306)/heroku_bace7cf727a523d")
	// LOCAL
	db, err := sql.Open("mysql", "isvacbanned:root@tcp(localhost:3306)/isvacbanned")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

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

	// Appr. 5 chars per num plus the comma.
	estimate := len(ids) * 6
	b := make([]byte, 0, estimate)

	for _, n := range ids {
		b = strconv.AppendInt(b, int64(n), 10)
		b = append(b, ',')
	}
	b = b[:len(b)-1]
	return string(b)
}
