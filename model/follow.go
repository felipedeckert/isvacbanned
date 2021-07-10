package model

import (
	"database/sql"
	"isvacbanned/util"
	"log"
	"strconv"

	//import mysql driver
	_ "github.com/go-sql-driver/mysql"
)

type FollowModel struct{}

type FollowModelInterface interface {
	FollowSteamUser(chatID int64, steamID, currNickname string, userID int64) int64
	UnfollowSteamUser(userID int64, steamID string) int64
	GetFollowerCountBySteamID(steamID string) (int64, error)
	GetAllIncompleteFollowedUsers() map[int64][]UsersFollowed
	GetUsersFollowed(userID int64) []UsersFollowed
	SetCurrNickname(userID int64, sanitizedActualNickname string)
	SetFollowedUserToCompleted(id []int64) int64
	GetUsersFollowedSummary(userID int64) map[bool]int
	IsFollowed(steamID string, userID int64) (string, int64, error)
}

var FollowModelClient FollowModelInterface = FollowModel{}

type UsersFollowed struct {
	ID           int64
	SteamID      string
	OldNickname  string
	CurrNickname string
	IsCompleted  bool
}

// FollowSteamUser links a telegram user to a steam user which is being followed
func (f FollowModel) FollowSteamUser(chatID int64, steamID, currNickname string, userID int64) int64 {

	stmt, err := util.GetDatabase().Prepare("INSERT INTO follow(chat_id, steam_id, user_id, old_nickname, curr_nickname) VALUES(?, ?, ?, ?, ?)")

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

// UnfollowSteamUser sets a followed player flag is_active to false
func (f FollowModel) UnfollowSteamUser(userID int64, steamID string) int64 {

	stmt, err := util.GetDatabase().Prepare("UPDATE follow SET is_active = false where user_id = ? AND steam_id = ?")

	if err != nil {
		panic(err.Error())
	}

	res, err := stmt.Exec(userID, steamID)

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
func (f FollowModel) GetFollowerCountBySteamID(steamID string) (int64, error) {

	row := util.GetDatabase().QueryRow(
		"SELECT COUNT(f.id) as count"+
			"	FROM follow f "+
			"	WHERE f.steam_id = ?", steamID)

	var count int64

	err := row.Scan(&count)

	return count, err
}

//GetAllIncompleteFollowedUsers get all fallowed steam user for every telegram user
func (f FollowModel) GetAllIncompleteFollowedUsers() map[int64][]UsersFollowed {

	rows, err := util.GetDatabase().Query("SELECT f.id, f.chat_id, f.steam_id, f.old_nickname, f.curr_nickname " +
		"FROM follow f " +
		"JOIN user u " +
		"ON f.user_id = u.id " +
		"WHERE f.is_completed <> true " +
		"AND f.is_active = true " +
		"AND u.is_active = true")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	m := make(map[int64][]UsersFollowed)
	for rows.Next() {
		var (
			id           int64
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
func (f FollowModel) GetUsersFollowed(userID int64) []UsersFollowed {

	rows, err := util.GetDatabase().Query("SELECT old_nickname, is_completed FROM follow WHERE user_id = ?", userID)

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

//SetCurrNickname updates de curr_nickname of a given player
func (f FollowModel) SetCurrNickname(userID int64, sanitizedActualNickname string) {

	stmt, err := util.GetDatabase().Prepare("UPDATE follow SET curr_nickname = ? where id = ?")

	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(sanitizedActualNickname, userID)

	defer stmt.Close()
	if err != nil {
		panic(err.Error())
	}

}

// SetFollowedUserToCompleted sets a player status to completed, and it will not be followed anymore
func (f FollowModel) SetFollowedUserToCompleted(id []int64) int64 {

	stmt, err := util.GetDatabase().Prepare("UPDATE follow SET is_completed = true where id in(?)")

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

// GetUsersFollowedSummary returns a summary of players followed, separated by ban status
func (f FollowModel) GetUsersFollowedSummary(userID int64) map[bool]int {

	rows, err := util.GetDatabase().Query(
		"SELECT COUNT(id) as count, is_completed"+
			"	FROM follow"+
			"	WHERE user_id = ?"+
			"	AND is_active = 1"+
			"	GROUP BY is_completed", userID)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	summary := make(map[bool]int, 0)
	for rows.Next() {
		var (
			count       int
			isCompleted bool
		)

		if err = rows.Scan(&count, &isCompleted); err != nil {
			if err != sql.ErrNoRows {
				panic(err.Error())
			}

			return nil
		}

		summary[isCompleted] = count
	}

	return summary
}

// IsFollowed checks if a user already follows a player
func (f FollowModel) IsFollowed(steamID string, userID int64) (string, int64, error) {

	row := util.GetDatabase().QueryRow("SELECT old_nickname, id FROM follow WHERE steam_id = ? AND user_id = ?", steamID, userID)

	var oldNickname string
	var id int64

	err := row.Scan(&oldNickname, &id)

	return oldNickname, id, err
}

func sliceToStringParam(ids []int64) string {
	if len(ids) == 0 {
		return ""
	}

	// Appr. 5 chars per num plus the comma.
	estimate := len(ids) * 6
	b := make([]byte, 0, estimate)

	for _, n := range ids {
		b = strconv.AppendInt(b, n, 10)
		b = append(b, ',')
	}
	b = b[:len(b)-1]
	return string(b)
}
