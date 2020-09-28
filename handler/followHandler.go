package handler

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"isvacbanned/messager"
	"isvacbanned/model"
	"isvacbanned/util"
)

type FollowHandler struct{}

var (
	FollowClient model.FollowModelClient
	MsgClient    MessageClient
)

func init() {
	UserModelClient = &model.UserModel{}
	FollowClient = &model.FollowModel{}
	MsgClient = &messager.MessageClient{}
}

//FollowHandler handles a follow request
func (f *FollowHandler) FollowHandler(m *tb.Message, bot *tb.Bot, steamID, currNickname string, userID int64) int64 {
	followersCount, err := FollowClient.GetFollowerCountBySteamID(steamID)

	UserModelClient.ActivateUser(userID)

	if err != nil {
		panic(err)
	}

	currNickname = util.SanitizeString(currNickname)

	dbID := FollowClient.FollowSteamUser(m.Chat.ID, steamID, currNickname, userID)

	message := getFollowResponse(currNickname, followersCount)

	MsgClient.SendMessage(bot, m.Sender, message)
	return dbID
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
