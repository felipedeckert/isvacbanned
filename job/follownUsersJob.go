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
	myMap := model.GetAllIncompletedFollowedUsers()

	idsToUpdate := make([]int, 0)

	for chatID, steamIDList := range myMap {
		for _, users := range steamIDList {
			player := service.UnmarshalPlayerByID(users.SteamID)
			playerData := player.Players[0]
			if playerData.VACBanned {
				sendMessageToUser(chatID, users.SteamID, playerData.DaysSinceLastBan)
				idsToUpdate = append(idsToUpdate, users.ID)
			}
		}
	}
	if len(idsToUpdate) > 0 {
		model.SetFollowedUserToCompleted(idsToUpdate)
	}
}

func sendMessageToUser(chatID int64, steamID string, daysSinceLastBan int) {
	token := "1324910657:AAFSlJn6TD9EeYNn35MEo-YphYlhYhqc_do"
	//token := os.Getenv("TOKEN")

	message := buildMessage(steamID, daysSinceLastBan)
	sendText := telegramAPIURL + token + telegramMethod + telegramChatIDParam + strconv.FormatInt(chatID, 10) + telegramTextParam + message
	http.Get(sendText)
}

func buildMessage(steamID string, daysSinceLastBan int) string {
	return "The user " + steamProfileURL + steamID + " has been BANNED " + strconv.Itoa(daysSinceLastBan) + " days ago!"
}
