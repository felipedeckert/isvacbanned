package handler

import (
	"database/sql"
	"isvacbanned/messenger"
	"isvacbanned/model"
	"isvacbanned/util"
	"log"

	tb "gopkg.in/tucnak/telebot.v2"
)

type FollowHandler struct{}

type FollowHandlerInterface interface {
	HandleFollowRequest(m *tb.Message, bot *tb.Bot, steamID, currNickname string, userID int64, isVACBanned bool) int64
}

var FollowHandlerClient FollowHandlerInterface = FollowHandler{}

//HandleFollowRequest handles a follow request
func (f FollowHandler) HandleFollowRequest(m *tb.Message, bot *tb.Bot, steamID, currNickname string, userID int64, isVACBanned bool) int64 {
	log.Printf("M=HandleFollowRequest steamID=%v\n", steamID)
	followersCount, err := model.FollowModelClient.GetFollowerCountBySteamID(steamID)

	log.Printf("M=HandleFollowRequest chatID=%v\n", m.Chat.ID)

	model.UserModelClient.ActivateUser(userID)

	if err != nil {
		panic(err)
	}

	var dbID int64

	followID, err := model.FollowModelClient.IsFollowed(steamID, userID)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	updateID := followID

	if followID == 0 {
		dbID = model.FollowModelClient.FollowSteamUser(m.Chat.ID, steamID, currNickname, userID)
		updateID = dbID
	}
	if isVACBanned {
		model.FollowModelClient.SetFollowedUserToCompleted([]int64{updateID})
	}

	response := util.GetFollowResponseMessage(currNickname, followersCount, followID, isVACBanned)

	messenger.MessengerClient.SendMessageToChat(bot, m.Chat, response)

	return dbID
}
