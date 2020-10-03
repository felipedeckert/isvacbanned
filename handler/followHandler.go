package handler

import (
	"database/sql"
	"fmt"
	"isvacbanned/messager"
	"isvacbanned/model"
	"log"

	tb "gopkg.in/tucnak/telebot.v2"
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
func (f *FollowHandler) FollowHandler(m *tb.Message, bot *tb.Bot, steamID, currNickname string, userID int64, isVACBanned bool) int64 {
	log.Printf("M=FollowHandler steamID=%v\n", steamID)
	followersCount, err := FollowClient.GetFollowerCountBySteamID(steamID)

	log.Printf("M=FollowHandler chatID=%v\n", m.Chat.ID)

	UserModelClient.ActivateUser(userID)

	if err != nil {
		panic(err)
	}

	var dbID int64

	followID, err := FollowClient.IsFollowed(steamID, userID)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	updateID := followID

	if followID == 0 {
		dbID = FollowClient.FollowSteamUser(m.Chat.ID, steamID, currNickname, userID)
		updateID = dbID
	}
	if isVACBanned {
		FollowClient.SetFollowedUserToCompleted([]int64{updateID})
	}

	message := getFollowResponse(currNickname, followersCount, followID, isVACBanned)

	MsgClient.SendMessageToChat(bot, m.Chat, message)
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
