package mock

import tb "gopkg.in/tucnak/telebot.v2"

type MsgClient struct {
	GetSendMessage       func(bot *tb.Bot, user *tb.User, message string)
	GetSendMessageToChat func(bot *tb.Bot, user *tb.Chat, message string)
}

func (m *MsgClient) SendMessage(bot *tb.Bot, user *tb.User, message string) {
	GetSendMessage(bot, user, message)
}

func (m *MsgClient) SendMessageToChat(bot *tb.Bot, user *tb.Chat, message string) {
	GetSendMessageToChat(bot, user, message)
}

var (
	GetSendMessage       func(bot *tb.Bot, user *tb.User, message string)
	GetSendMessageToChat func(bot *tb.Bot, user *tb.Chat, message string)
)
