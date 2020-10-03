package handler

import (
	"fmt"
	"isvacbanned/model"
	"log"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

const maxMessageLength int = 3800

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
	total := len(followedUsers)

	if total == 0 {
		return "You're not following any player yet!"
	}

	msg, count := getPlayersAndStatusAsShoppingList(followedUsers)

	return fmt.Sprintf("You're following these users: \n%v And %v more...", msg, total-count)
}

func getPlayersAndStatusAsShoppingList(followedUsers []model.UsersFollowed) (string, int) {
	var str strings.Builder
	var count int
	for _, user := range followedUsers {
		count++
		status := "NOT BANNED"
		if user.IsCompleted {
			status = "BANNED"
		}
		str.WriteString(user.OldNickname + " : " + status + ",\n")
		if str.Len() > maxMessageLength {
			break
		}
	}

	return str.String(), count
}
