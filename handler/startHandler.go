package handler

import (
	"isvacbanned/util"
	"log"

	tb "gopkg.in/tucnak/telebot.v2"
)

//StartHandler handles show requests
func StartHandler(m *tb.Message, bot *tb.Bot) {
	message := util.GetStartResponse(m.Sender.FirstName)

	_, err := bot.Send(m.Chat, message)
	if err != nil {
		log.Printf("M=StartHandler err=%s", err.Error())
	}

}
