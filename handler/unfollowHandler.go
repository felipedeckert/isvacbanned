package handler

import (
	"isvacbanned/messenger"
	"isvacbanned/model"
	"isvacbanned/util"

	tb "gopkg.in/tucnak/telebot.v2"
)

type UnfollowHandler struct{}

type UnfollowHandlerInterface interface {
	HandleUnfollowRequest(m *tb.Message, bot *tb.Bot, steamID string, userID int64)
}

var UnfollowHandlerClient UnfollowHandlerInterface = UnfollowHandler{}

//HandleUnfollowRequest handles a follow request
func (f UnfollowHandler) HandleUnfollowRequest(m *tb.Message, bot *tb.Bot, steamID string, userID int64) {

	rows := model.FollowModelClient.UnfollowSteamUser(userID, steamID)

	var message string

	if rows != 1 {
		message = getUnsuccessfulUnfollowResponse()
	} else {
		message = getSuccessfulUnfollowResponse(steamID)
	}

	messenger.MessengerClient.SendMessageToChat(bot, m.Chat, message)

}

func getUnsuccessfulUnfollowResponse() string {
	return "Unable to unfollow this player. Are you sure you follow him/her?"
}

func getSuccessfulUnfollowResponse(steamID string) string {
	return "You will NOT receive more updates about this player: " + util.SteamProfileURL + steamID
}
