package handler

import (
	"fmt"
	"isvacbanned/model"
	"log"
	"regexp"

	tb "gopkg.in/tucnak/telebot.v2"
)

func FollowHandler(m *tb.Message, bot *tb.Bot, steamID, currNickname string, userID int64) {

	followersCount, err := model.GetFollowerCountBySteamID(steamID)

	if err != nil {
		panic(err)
	}

	currNickname = SanitizeString(currNickname)

	model.FollowSteamUser(m.Chat.ID, steamID, currNickname, userID)

	message := getFollowResponse(currNickname, followersCount)

	bot.Send(m.Sender, message)
}

func SanitizeString(input string) string {
	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		log.Fatal(err)
	}
	sanitizedInput := re.ReplaceAllString(input, "")

	if len(sanitizedInput) == 0 {
		sanitizedInput = "Unreadable nickname"
	}

	return sanitizedInput
}

func getFollowResponse(currNickname string, followersCount int64) string {

	message := fmt.Sprintf("Following player %v, ", currNickname)

	if followersCount > 0 {
		message += fmt.Sprintf("which is being followes by %v other users.", followersCount)
	} else {
		message += "you're the first to follow this player!"
	}

	return message
}
