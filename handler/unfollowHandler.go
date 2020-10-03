package handler

import (
	"isvacbanned/messager"
	"isvacbanned/model"
	"isvacbanned/util"

	tb "gopkg.in/tucnak/telebot.v2"
)

type UnfollowHandler struct{}

func init() {
	FollowClient = &model.FollowModel{}
	MsgClient = &messager.MessageClient{}
}

//UnfollowHandler handles a follow request
func (f *UnfollowHandler) UnfollowHandler(m *tb.Message, bot *tb.Bot, steamID string, userID int64) {

	rows := FollowClient.UnfollowSteamUser(userID, steamID)

	var message string

	if rows != 1 {
		message = getUnsuccessfulUnfollowResponse()
	} else {
		message = getSuccessfulUnfollowResponse(steamID)
	}

	MsgClient.SendMessageToChat(bot, m.Chat, message)

}

func getUnsuccessfulUnfollowResponse() string {
	return "Unable to unfollow this player. Are you sure you follow him/her?"
}

func getSuccessfulUnfollowResponse(steamID string) string {
	return "You will NOT receive more updates about this player: " + util.SteamProfileURL + steamID
}
