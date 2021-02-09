package service

import (
	"database/sql"
	"isvacbanned/model"
	"log"

	tb "gopkg.in/tucnak/telebot.v2"
)

type userService struct{}

type UserServiceInterface interface {
	getUserID(chat *tb.Chat) int64
}

var UserServiceClient UserServiceInterface = userService{}

func (u userService) getUserID(chat *tb.Chat) int64 {
	log.Printf("M=getUserID telegramID=%v\n", chat.ID)
	id, err := model.UserModelClient.GetUserID(chat.ID)
	if err != nil && err == sql.ErrNoRows {
		id = model.UserModelClient.CreateUser(chat.FirstName, chat.Username, chat.ID)
	}

	log.Printf("M=getUserID userID=%v\n", id)
	return id
}
