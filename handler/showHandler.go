package handler

import (
	"isvacbanned/model"
	"log"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

var followClient model.FollowModelClient

func init() {
	followClient = &model.FollowModel{}
}

//ShowHandler handles show requests
func ShowHandler(m *tb.Message, bot *tb.Bot, userID int64) {
	followedUsers := followClient.GetUsersFollowed(userID)

	log.Printf("M=ShowHandler L=I userID=%v usersFollowedCount=%v\n", m.Chat.ID, len(followedUsers))

	message := getShowResponse(followedUsers)

	_, err := bot.Send(m.Chat, message)

	if err != nil {
		log.Printf("M=ShowHandler L=E userID=%v err=%v\n", m.Chat.ID, err.Error())
	}
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
