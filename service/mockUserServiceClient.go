package service

import tb "gopkg.in/tucnak/telebot.v2"

type UserServiceMock struct{}

func (u UserServiceMock) getUserID(chat *tb.Chat) int64 {

	return int64(321)
}
