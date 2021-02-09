package handler

import tb "gopkg.in/tucnak/telebot.v2"

type MessageClient interface {
	SendMessage(bot *tb.Bot, user *tb.User, message string)
	SendMessageToChat(bot *tb.Bot, user *tb.Chat, message string)
}
