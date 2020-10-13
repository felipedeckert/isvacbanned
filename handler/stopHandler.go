package handler

import (
	"isvacbanned/model"
	"isvacbanned/util"

	tb "gopkg.in/tucnak/telebot.v2"
)

//StopHandler handles show requests
func StopHandler(m *tb.Message, bot *tb.Bot, userID int64) {
	model.UserModelClient.InactivateUser(userID)

	message := util.GetStopResponse()

	bot.Send(m.Chat, message)
}
