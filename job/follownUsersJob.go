package job

import (
	"fmt"
	"isvacbanned/model"
	"isvacbanned/service"
	"net/http"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"
)

const steamProfileURL = "https://steamcommunity.com/profiles/"
const telegramAPIURL = "https://api.telegram.org/bot"
const telegramMethod = "/sendMessage"
const telegramChatIDParam = "?chat_id="
const telegramTextParam = "&text="

// RunScheduler sets up the scheduler to check users status
func RunScheduler() {
	fmt.Println("RunScheduler")
	scheduler := gocron.NewScheduler(time.UTC)

	scheduler.Every(30).Seconds().Do(checkFollownUsers)

	scheduler.StartAsync()
}

func checkFollownUsers() {
	fmt.Println("checkFollownUsers")
	usersIncompleted := model.GetAllIncompletedFollowedUsers()
	var idsToUpdate []int
	for chatID, steamIDList := range usersIncompleted {
		for _, users := range steamIDList {
			idsToUpdate = handleFollowedUser(users, chatID)
		}
	}
	if len(idsToUpdate) > 0 {
		model.SetFollowedUserToCompleted(idsToUpdate)
	}
}

func handleFollowedUser(user model.UsersFollowed, chatID int64) []int {
	idsToUpdate := make([]int, 0)
	player := service.UnmarshalPlayerByID(user.SteamID)
	playerData := player.Players[0]

	actualNickname := service.GetPlayerCurrentNickname(user.SteamID)

	if playerData.VACBanned {
		sendMessageToUser(buildBanMessage(user.OldNickname, actualNickname, user.SteamID, playerData.DaysSinceLastBan), chatID)
		idsToUpdate = append(idsToUpdate, user.ID)
	} else if user.CurrNickname != actualNickname {
		sendMessageToUser(buildNicknameChangedMessage(user.OldNickname, user.CurrNickname, user.SteamID), chatID)
	}
	return idsToUpdate
}

func sendMessageToUser(message string, chatID int64) {
	token := "1324910657:AAFSlJn6TD9EeYNn35MEo-YphYlhYhqc_do"
	//token := os.Getenv("TOKEN")

	sendMessageURL := telegramAPIURL + token + telegramMethod + telegramChatIDParam + strconv.FormatInt(chatID, 10) + telegramTextParam + message
	http.Get(sendMessageURL)
}

func buildNicknameChangedMessage(oldNickname, currNickname, steamID string) string {
	return "The user you followed as " + oldNickname + " Steam Profile: " + steamProfileURL + steamID + " is now under the nickname " + currNickname
}

func buildBanMessage(oldNickname, currNickname, steamID string, daysSinceLastBan int) string {
	// if the player hanst changed nickname no reason to return redundant message
	changedNickPhrase := ""
	if oldNickname != currNickname {
		changedNickPhrase = ", now under the nickname " + currNickname + ","
	}
	return "The user you followed as " + oldNickname + changedNickPhrase + " Steam Profile: " + steamProfileURL + steamID + ", has been BANNED " + strconv.Itoa(daysSinceLastBan) + " days ago!"
}
