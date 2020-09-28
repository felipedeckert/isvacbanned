package mock

import tb "gopkg.in/tucnak/telebot.v2"

type MsgClient struct {
	GetSendMessage func(bot *tb.Bot, user *tb.User, message string)
}

func (m *MsgClient) SendMessage(bot *tb.Bot, user *tb.User, message string) {
	GetSendMessage(bot, user, message)
}

var (
	GetSendMessage func(bot *tb.Bot, user *tb.User, message string)
)
