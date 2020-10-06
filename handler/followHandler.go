package handler

import (
	"database/sql"
	"fmt"
	"isvacbanned/messenger"
	"isvacbanned/model"
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

	response := getFollowResponse(currNickname, followersCount, followID, isVACBanned)

	messenger.MessengerClient.SendMessageToChat(bot, m.Chat, response)

	return dbID
}

func getFollowResponse(currNickname string, followersCount, followID int64, isVACBanned bool) string {

	var status string = "NOT banned (yet)."
	if isVACBanned {
		status = "BANNED (yay)."
	}

	if followID > 0 {
		return fmt.Sprintf("You already follow this player. Its current status is: %v", status)
	}

	message := fmt.Sprintf("Following player %v, status=%v", currNickname, status)

	if !isVACBanned {
		if followersCount > 0 {
			message += fmt.Sprintf(" Which is being followed by %v other users.", followersCount)
		} else {
			message += " You're the first to follow this player!"
		}
	} else {
		message += " You will NOT receive updates about this player since it's banned!"
	}
	return message
}
