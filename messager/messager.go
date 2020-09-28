package messager

import tb "gopkg.in/tucnak/telebot.v2"

type MessageClient struct{}

func (m *MessageClient) SendMessage(bot *tb.Bot, user *tb.User, message string) {
	bot.Send(user, message)
}
