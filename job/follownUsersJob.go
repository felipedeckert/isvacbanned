package job

import (
	"isvacbanned/messenger"
	"isvacbanned/model"
	"isvacbanned/service"
	"isvacbanned/util"
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

// RunScheduler sets up the scheduler to check users status
func RunScheduler() {
	log.Printf("M=RunScheduler\n")
	scheduler := gocron.NewScheduler(time.UTC)

	scheduler.Every(60).Minutes().Do(checkFollownUsersBan)
	scheduler.Every(10).Minutes().Do(checkFollownUsersNickname)

	scheduler.StartAsync()
}

func checkFollownUsersBan() {
	usersIncompleted := model.FollowModelClient.GetAllIncompletedFollowedUsers()
	var usersToComplete []int64
	for chatID, steamIDList := range usersIncompleted {
		for _, users := range steamIDList {
			usersToComplete = append(usersToComplete, validateBanStatusAndSendMessage(users, chatID)...)
		}
	}
	if len(usersToComplete) > 0 {
		//Once a player status is set to completed, this player will not be returned in the GetAllIncompletedFollowedUsers query
		model.FollowModelClient.SetFollowedUserToCompleted(usersToComplete)
	}

	log.Printf("M=checkFollownUsersBan usersToCompleteCount=%v\n", len(usersToComplete))
}

func checkFollownUsersNickname() {
	usersIncompleted := model.FollowModelClient.GetAllIncompletedFollowedUsers()

	for chatID, steamIDList := range usersIncompleted {
		for _, users := range steamIDList {
			validateNicknameAndSendMessage(users, chatID)
		}
	}
}

func validateBanStatusAndSendMessage(user model.UsersFollowed, chatID int64) []int64 {
	idsToUpdate := make([]int64, 0)
	player := service.PlayerServiceClient.GetPlayerStatus(user.SteamID)
	if len(player.Players) > 0 {

		playerData := player.Players[0]
		if playerData.VACBanned {
			actualNickname := service.PlayerServiceClient.GetPlayerCurrentNickname(user.SteamID)
			if actualNickname != "" {
				log.Printf("M=validateBanStatusAndSendMessage steamID=%v status=banned\n", user.SteamID)
				messenger.MessengerClient.SendMessageToUser(util.GetBanMessage(user.OldNickname, actualNickname, user.SteamID, playerData.DaysSinceLastBan), chatID)
				idsToUpdate = append(idsToUpdate, user.ID)
			}
		}
	}
	return idsToUpdate
}

func validateNicknameAndSendMessage(user model.UsersFollowed, chatID int64) {
	actualNickname := service.PlayerServiceClient.GetPlayerCurrentNickname(user.SteamID)

	if actualNickname == "" {
		log.Printf("M=validateNicknameAndSendMessage L=E steamID=%v status=emptyCurrentNickname\n", user.SteamID)
	} else if user.CurrNickname != actualNickname {
		log.Printf("M=validateNicknameAndSendMessage steamID=%v status=changedNickname\n", user.SteamID)
		model.FollowModelClient.SetCurrNickname(user.ID, actualNickname)
		messenger.MessengerClient.SendMessageToUser(util.GetNicknameChangedMessage(user.OldNickname, user.CurrNickname, actualNickname, user.SteamID), chatID)
	}
}
