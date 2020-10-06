package job

import (
	"isvacbanned/messenger"
	"isvacbanned/model"
	"isvacbanned/service"
	"isvacbanned/util"
	"log"
	"strconv"
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
	playerData := player.Players[0]

	if playerData.VACBanned {
		actualNickname := service.PlayerServiceClient.GetPlayerCurrentNickname(user.SteamID)
		log.Printf("M=validateBanStatusAndSendMessage steamID=%v status=banned\n", user.SteamID)
		messenger.MessengerClient.SendMessageToUser(buildBanMessage(user.OldNickname, actualNickname, user.SteamID, playerData.DaysSinceLastBan), chatID)
		idsToUpdate = append(idsToUpdate, user.ID)
	}
	return idsToUpdate
}

func validateNicknameAndSendMessage(user model.UsersFollowed, chatID int64) {
	actualNickname := service.PlayerServiceClient.GetPlayerCurrentNickname(user.SteamID)

	if user.CurrNickname != actualNickname {
		log.Printf("M=validateNicknameAndSendMessage steamID=%v status=changedNickname\n", user.SteamID)
		model.FollowModelClient.SetCurrNickname(user.ID, actualNickname)
		messenger.MessengerClient.SendMessageToUser(buildNicknameChangedMessage(user.OldNickname, user.CurrNickname, actualNickname, user.SteamID), chatID)
	}
}

func buildNicknameChangedMessage(oldNickname, recentNickname, currNickname, steamID string) string {
	diffRecentNickname := ""
	if oldNickname != recentNickname {
		diffRecentNickname = ", recently playing as \"" + recentNickname + "\""
	}

	return "The user you followed as \"" + oldNickname + "\"" + diffRecentNickname + ", Steam Profile: " + util.SteamProfileURL + steamID + ", is now under the nickname \"" + currNickname + "\""
}

func buildBanMessage(oldNickname, currNickname, steamID string, daysSinceLastBan int) string {
	// if the player hasn't changed nickname no reason to return redundant message
	changedNickPhrase := ""
	if oldNickname != currNickname {
		changedNickPhrase = ", now under the nickname " + currNickname
	}
	return "The user you followed as " + oldNickname + changedNickPhrase + ", Steam Profile: " + util.SteamProfileURL + steamID + ", has been BANNED " + strconv.Itoa(daysSinceLastBan) + " days ago! You won't be notified about this player anymore."
}
