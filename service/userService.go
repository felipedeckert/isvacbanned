package service

import (
	"database/sql"
	"isvacbanned/model"

	tb "gopkg.in/tucnak/telebot.v2"
)

func getUserID(user *tb.User) int64 {

	id, err := model.GetUserID(user.ID)
	if err != nil && err == sql.ErrNoRows {
		id = model.CreateUser(user.FirstName, user.Username, user.ID)
	}

	return id
}
