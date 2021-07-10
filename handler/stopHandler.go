package handler

import (
	"isvacbanned/model"
	"isvacbanned/util"
	"log"

	tb "gopkg.in/tucnak/telebot.v2"
)

//StopHandler handles show requests
func StopHandler(m *tb.Message, bot *tb.Bot, userID int64) {
	model.UserModelClient.InactivateUser(userID)

	message := util.GetStopResponse()

	_, err := bot.Send(m.Chat, message)
	if err != nil {
		log.Printf("M=StopHandler err=%s", err.Error())
	}
}
