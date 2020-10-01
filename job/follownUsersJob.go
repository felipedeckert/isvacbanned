package job

import (
	"isvacbanned/model"
	"isvacbanned/service"
	"isvacbanned/util"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"
)

const telegramAPIURL = "https://api.telegram.org/bot"
const telegramMethod = "/sendMessage"
const telegramChatIDParam = "?chat_id="
const telegramTextParam = "&text="

var followModelClient *model.FollowModel

func init() {
	followModelClient = &model.FollowModel{}
}

// RunScheduler sets up the scheduler to check users status
func RunScheduler() {
	log.Printf("M=RunScheduler\n")
	scheduler := gocron.NewScheduler(time.UTC)

	scheduler.Every(60).Minutes().Do(checkFollownUsersBan)
	scheduler.Every(10).Minutes().Do(checkFollownUsersNickname)

	scheduler.StartAsync()
}

func checkFollownUsersBan() {
	usersIncompleted := followModelClient.GetAllIncompletedFollowedUsers()
	var usersToComplete []int64
	for chatID, steamIDList := range usersIncompleted {
		for _, users := range steamIDList {
			usersToComplete = append(usersToComplete, hasPlayerBeenBanned(users, chatID)...)
		}
	}
	if len(usersToComplete) > 0 {
		//Once a player status is set to completed, this player will not be returned in the GetAllIncompletedFollowedUsers query
		followModelClient.SetFollowedUserToCompleted(usersToComplete)
	}

	log.Printf("M=checkFollownUsersBan usersToCompleteCount=%v\n", len(usersToComplete))
}

func checkFollownUsersNickname() {
	usersIncompleted := followModelClient.GetAllIncompletedFollowedUsers()

	for chatID, steamIDList := range usersIncompleted {
		for _, users := range steamIDList {
			hasPlayerChangedNickname(users, chatID)
		}
	}
}

func hasPlayerBeenBanned(user model.UsersFollowed, chatID int64) []int64 {
	idsToUpdate := make([]int64, 0)
	player := service.GetPlayerStatus(user.SteamID)
	playerData := player.Players[0]

	if playerData.VACBanned {
		actualNickname := service.GetPlayerCurrentNickname(user.SteamID)
		log.Printf("M=hasPlayerBeenBanned steamID=%v status=banned\n", user.SteamID)
		sendMessageToUser(buildBanMessage(user.OldNickname, actualNickname, user.SteamID, playerData.DaysSinceLastBan), chatID)
		idsToUpdate = append(idsToUpdate, user.ID)
	}
	return idsToUpdate
}

func hasPlayerChangedNickname(user model.UsersFollowed, chatID int64) {
	actualNickname := service.GetPlayerCurrentNickname(user.SteamID)

	if user.CurrNickname != actualNickname {
		log.Printf("M=hasPlayerBeenBanned steamID=%v status=changedNickname\n", user.SteamID)
		followModelClient.SetCurrNickname(user.ID, actualNickname)
		sendMessageToUser(buildNicknameChangedMessage(user.OldNickname, actualNickname, user.SteamID), chatID)
	}
}

func sendMessageToUser(message string, chatID int64) {
	var token string

	if util.LOCAL {
		token = "1324910657:AAFSlJn6TD9EeYNn35MEo-YphYlhYhqc_do"
	} else {
		token = os.Getenv("TOKEN")
	}

	sendMessageURL := telegramAPIURL + token + telegramMethod + telegramChatIDParam + strconv.FormatInt(chatID, 10) + telegramTextParam + message
	_, err := http.Get(sendMessageURL)
	if err != nil {
		panic(err)
	}
}

func buildNicknameChangedMessage(oldNickname, currNickname, steamID string) string {
	return "The user you followed as " + oldNickname + " Steam Profile: " + util.SteamProfileURL + steamID + " is now under the nickname " + currNickname
}

func buildBanMessage(oldNickname, currNickname, steamID string, daysSinceLastBan int) string {
	// if the player hasn't changed nickname no reason to return redundant message
	changedNickPhrase := ""
	if oldNickname != currNickname {
		changedNickPhrase = ", now under the nickname " + currNickname
	}
	return "The user you followed as " + oldNickname + changedNickPhrase + ", Steam Profile: " + util.SteamProfileURL + steamID + ", has been BANNED " + strconv.Itoa(daysSinceLastBan) + " days ago! You won't be notified about this player anymore."
}
