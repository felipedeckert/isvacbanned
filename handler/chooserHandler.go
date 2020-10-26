package handler

import (
	"isvacbanned/util"
	"math/rand"

	tb "gopkg.in/tucnak/telebot.v2"
)

//HandleChooserRequest handles show requests
func HandleChooserRequest(m *tb.Message, bot *tb.Bot, userID int64) {

	res := rand.Intn(2)

	bot.Send(m.Chat, util.GetChooserResponse(res))
}
