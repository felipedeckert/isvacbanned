package handler

import (
	"isvacbanned/util"

	tb "gopkg.in/tucnak/telebot.v2"
)

//StartHandler handles show requests
func StartHandler(m *tb.Message, bot *tb.Bot) {
	message := util.GetStartResponse(m.Sender.FirstName)

	bot.Send(m.Chat, message)
}
