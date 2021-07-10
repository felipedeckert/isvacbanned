package handler

import (
	"isvacbanned/util"
	"log"
	"math/rand"

	tb "gopkg.in/tucnak/telebot.v2"
)

//HandleChooserRequest handles show requests
func HandleChooserRequest(m *tb.Message, bot *tb.Bot) {

	res := rand.Intn(2)

	_, err := bot.Send(m.Chat, util.GetChooserResponse(res))
	if err != nil {
		log.Printf("M=HandleChooserRequest err=%s", err.Error())
	}
}
