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

	scheduler.Every(30).Seconds().Do(checkFollownUsers)

	scheduler.StartAsync()
}

func checkFollownUsers() {
	usersIncompleted := followModelClient.GetAllIncompletedFollowedUsers()
	var usersToComplete []int
	for chatID, steamIDList := range usersIncompleted {
		for _, users := range steamIDList {
			usersToComplete = handleFollowedUser(users, chatID)
		}
	}
	if len(usersToComplete) > 0 {
		//Once a player status is set to completed, this player will not be returned in the GetAllIncompletedFollowedUsers query
		followModelClient.SetFollowedUserToCompleted(usersToComplete)
	}

	log.Printf("M=checkFollownUsers usersToCompleteCount=%v\n", len(usersToComplete))
}

func handleFollowedUser(user model.UsersFollowed, chatID int64) []int {
	idsToUpdate := make([]int, 0)
	player := service.GetPlayerStatus(user.SteamID)
	playerData := player.Players[0]

	actualNickname := service.GetPlayerCurrentNickname(user.SteamID)

	sanitizedActualNickname := util.SanitizeString(actualNickname)

	if playerData.VACBanned {
		log.Printf("M=handleFollowedUser steamID=%v status=banned\n", user.SteamID)
		sendMessageToUser(buildBanMessage(user.OldNickname, actualNickname, user.SteamID, playerData.DaysSinceLastBan), chatID)
		idsToUpdate = append(idsToUpdate, user.ID)
	} else if user.CurrNickname != sanitizedActualNickname {
		log.Printf("M=handleFollowedUser steamID=%v status=changedNickname\n", user.SteamID)
		followModelClient.SetCurrNickname(user.ID, sanitizedActualNickname)
		sendMessageToUser(buildNicknameChangedMessage(user.OldNickname, actualNickname, user.SteamID), chatID)
	}
	return idsToUpdate
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
	if oldNickname != util.SanitizeString(currNickname) {
		changedNickPhrase = ", now under the nickname " + currNickname
	}
	return "The user you followed as " + oldNickname + changedNickPhrase + ", Steam Profile: " + util.SteamProfileURL + steamID + ", has been BANNED " + strconv.Itoa(daysSinceLastBan) + " days ago! You won't be notified about this player anymore."
}
