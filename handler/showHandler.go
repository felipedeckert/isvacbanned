package handler

import (
	"isvacbanned/model"
	"log"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

const maxMessageLength int = 4000

var followClient model.FollowModelClient

func init() {
	followClient = &model.FollowModel{}
}

//ShowHandler handles show requests
func ShowHandler(m *tb.Message, bot *tb.Bot, userID int64) {
	followedUsers := followClient.GetUsersFollowed(userID)

	log.Printf("M=ShowHandler L=I userID=%v usersFollowedCount=%v\n", m.Chat.ID, len(followedUsers))

	sendShowResponse(followedUsers, bot, m)
}

func sendShowResponse(followedUsers []model.UsersFollowed, bot *tb.Bot, m *tb.Message) {
	total := len(followedUsers)

	if total == 0 {
		bot.Send(m.Chat, "You're not following any player yet!")
	}

	sendMessageBatch(followedUsers, bot, m)
}

func sendMessageBatch(followedUsers []model.UsersFollowed, bot *tb.Bot, m *tb.Message) {
	var sb strings.Builder
	var count int
	var status string
	var prefix string
	sb.WriteString("You're following these users: \n")
	for _, user := range followedUsers {
		count++
		prefix = ""
		status = "NOT BANNED"
		if user.IsCompleted {
			prefix = "❌ "
			status = "BANNED ❌"
		}
		sb.WriteString(prefix + user.OldNickname + " : " + status + ",\n")
		if sb.Len() > maxMessageLength {
			sendAndResetBuffer(&sb, bot, m)
		}
	}

	if sb.Len() > 0 {
		bot.Send(m.Chat, sb.String())
	}

}

func sendAndResetBuffer(sb *strings.Builder, bot *tb.Bot, m *tb.Message) {
	_, err := bot.Send(m.Chat, sb.String())

	sb.Reset()

	if err != nil {
		log.Printf("M=ShowHandler L=E userID=%v err=%v\n", m.Chat.ID, err.Error())
	}
}
