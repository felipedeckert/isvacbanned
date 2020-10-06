package handler

import (
	"fmt"
	"isvacbanned/model"
	"log"

	tb "gopkg.in/tucnak/telebot.v2"
)

//HandleSummaryRequest handles show requests
func HandleSummaryRequest(m *tb.Message, bot *tb.Bot, userID int64) {
	summary := model.FollowModelClient.GetUsersFollowedSummary(userID)

	log.Printf("M=HandleSummaryRequest L=I userID=%v \n", m.Chat.ID)

	sendSummaryResponse(summary, bot, m)
}

func sendSummaryResponse(summary map[bool]int, bot *tb.Bot, m *tb.Message) {

	messageEnd := ""

	if summary[false]+summary[true] == 0 {
		messageEnd = ", let's start tracking some suspects!"
	} else if summary[true] == 0 {
		messageEnd = fmt.Sprintf(", no one have been banned yet. Don't worry VAC will get them!")
	} else {
		messageEnd = fmt.Sprintf(", of witch %v have been banned, a performance of %.2f%%. Keep up the good work!", summary[true], (float64(summary[true])/float64(summary[false]+summary[true]))*100)
	}

	message := fmt.Sprintf("You follow %v players%v", summary[true]+summary[false], messageEnd)

	bot.Send(m.Chat, message)
}
