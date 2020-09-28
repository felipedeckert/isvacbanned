package handler

import (
	"fmt"
	"isvacbanned/messager"
	"isvacbanned/model"
	"log"
	"regexp"

	tb "gopkg.in/tucnak/telebot.v2"
)

type FollowHandler struct{}

var (
	FollowClient FollowModelClient
	MsgClient    MessageClient
)

func init() {
	FollowClient = &model.FollowModel{}
	MsgClient = &messager.MessageClient{}
}

//FollowHandler handles a follow request
func (f *FollowHandler) FollowHandler(m *tb.Message, bot *tb.Bot, steamID, currNickname string, userID int64) int64 {
	fmt.Println("AAAAAAAAAAAAAAAAAA")
	followersCount, err := FollowClient.GetFollowerCountBySteamID(steamID)

	if err != nil {
		panic(err)
	}

	currNickname = SanitizeString(currNickname)

	dbID := FollowClient.FollowSteamUser(m.Chat.ID, steamID, currNickname, userID)

	message := getFollowResponse(currNickname, followersCount)

	MsgClient.SendMessage(bot, m.Sender, message)
	return dbID
}

// SanitizeString removes all non acsii chars form a string
func SanitizeString(input string) string {
	re, err := regexp.Compile(`[^\x00-\x7F]`)
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
