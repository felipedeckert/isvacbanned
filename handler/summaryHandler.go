package handler

import (
	"isvacbanned/model"
	"isvacbanned/util"
	"log"

	tb "gopkg.in/tucnak/telebot.v2"
)

//HandleSummaryRequest handles show requests
func HandleSummaryRequest(m *tb.Message, bot *tb.Bot, userID int64) {
	summary := model.FollowModelClient.GetUsersFollowedSummary(userID)

	log.Printf("M=HandleSummaryRequest L=I userID=%v \n", m.Chat.ID)

	bot.Send(m.Chat, util.GetSummaryResponse(summary))
}
