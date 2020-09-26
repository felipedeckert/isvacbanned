package handler

import (
	"isvacbanned/model"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

func ShowHandler(m *tb.Message, bot *tb.Bot, userID int64) {
	followedUsers := model.GetUsersFollowed(userID)

	message := getShowResponse(followedUsers)

	bot.Send(m.Sender, message)
}

func getShowResponse(followedUsers []model.UsersFollowed) string {
	if len(followedUsers) == 0 {
		return "You're not following any player yet!"
	}

	return "You're following these users: \n" + getPlayersAndStatusAsShoppingList(followedUsers)
}

func getPlayersAndStatusAsShoppingList(followedUsers []model.UsersFollowed) string {
	var str strings.Builder

	for _, user := range followedUsers {
		status := "NOT BANNED"
		if user.IsCompleted {
			status = "BANNED"
		}
		str.WriteString(user.OldNickname + " : " + status + ",\n")
	}

	return str.String()
}
