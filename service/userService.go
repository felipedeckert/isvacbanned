package service

import (
	"database/sql"
	"isvacbanned/model"
	"log"

	tb "gopkg.in/tucnak/telebot.v2"
)

var userModelClient model.UserModelClient

func init() {
	userModelClient = &model.UserModel{}
}

func getUserID(user *tb.User) int64 {
	log.Printf("M=getUserID telegramID=%v\n", user.ID)
	id, err := userModelClient.GetUserID(user.ID)
	if err != nil && err == sql.ErrNoRows {
		id = userModelClient.CreateUser(user.FirstName, user.Username, user.ID)
	}

	return id
}
