package handler

import (
	"isvacbanned/model"

	tb "gopkg.in/tucnak/telebot.v2"
)

//StopHandler handles show requests
func StopHandler(m *tb.Message, bot *tb.Bot, userID int64) {
	model.UserModelClient.InactivateUser(userID)

	message := getStopResponse()

	bot.Send(m.Chat, message)
}

func getStopResponse() string {
	return "You will not be notified about any player anymore! Follow another player to start receiving news about all the players you followed."
}
