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

func getUserID(chat *tb.Chat) int64 {
	log.Printf("M=getUserID telegramID=%v\n", chat.ID)
	id, err := userModelClient.GetUserID(chat.ID)
	if err != nil && err == sql.ErrNoRows {
		id = userModelClient.CreateUser(chat.FirstName, chat.Username, chat.ID)
	}

	return id
}
