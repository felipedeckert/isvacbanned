package service

import (
	"database/sql"
	"fmt"
	"isvacbanned/model"

	tb "gopkg.in/tucnak/telebot.v2"
)

var userModelClient UserModelClient

func init() {
	userModelClient = &model.UserModel{}
}

func getUserID(user *tb.User) int64 {

	id, err := userModelClient.GetUserID(user.ID)
	if err != nil && err == sql.ErrNoRows {
		id = userModelClient.CreateUser(user.FirstName, user.Username, user.ID)
	}
	fmt.Println(id)
	return id
}
