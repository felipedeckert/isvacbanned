package job

import (
	"isvacbanned/model"
	"isvacbanned/service"
	"net/http"
	"os"
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
	scheduler := gocron.NewScheduler(time.UTC)

	scheduler.Every(2).Minutes().Do(checkFollownUsers)
}

func checkFollownUsers() {
	myMap := model.GetAllFollownUsers()

	for chatID, steamIDList := range myMap {
		for _, steamID := range steamIDList {
			player := service.UnmarshalPlayerByID(steamID)
			playerData := player.Players[0]
			if playerData.VACBanned {
				sendMessageToUser(chatID, steamID, playerData.DaysSinceLastBan)
			}
		}
	}
}

func sendMessageToUser(chatID int64, steamID string, daysSinceLastBan int) {
	token := os.Getenv("TOKEN")

	message := buildMessage(steamID)
	sendText := telegramAPIURL + token + telegramMethod + telegramChatIDParam + strconv.FormatInt(chatID, 10) + telegramTextParam + message
	http.Get(sendText)
}

func buildMessage(steamID string) string {
	return "The user " + steamProfileURL + steamID + " has been BANNED!"
}
