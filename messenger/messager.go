package messenger

import (
	"net/http"
	"os"
	"strconv"

	tb "gopkg.in/tucnak/telebot.v2"
)

const telegramAPIURL = "https://api.telegram.org/bot"
const telegramMethod = "/sendMessage"
const telegramChatIDParam = "?chat_id="
const telegramTextParam = "&text="

//Messenger is the real implementation of the Messenger
type Messenger struct{}

//MessengerInterface is the interface that defines the behavior of the Messenger
type MessengerInterface interface {
	SendMessage(bot *tb.Bot, user *tb.User, message string)
	SendMessageToChat(bot *tb.Bot, chat *tb.Chat, message string)
	SendMessageToUser(message string, chatID int64)
}

//MessengerClient client to send messages
var MessengerClient MessengerInterface = Messenger{}

//SendMessage sends a message to user via bot
func (m Messenger) SendMessage(bot *tb.Bot, user *tb.User, message string) {
	bot.Send(user, message)
}

//SendMessageToChat sends a message to chat via bot
func (m Messenger) SendMessageToChat(bot *tb.Bot, chat *tb.Chat, message string) {
	bot.Send(chat, message)
}

//SendMessageToUser sends a message to user/group via telegram API
func (m Messenger) SendMessageToUser(message string, chatID int64) {
	token := os.Getenv("TOKEN")

	sendMessageURL := telegramAPIURL + token + telegramMethod + telegramChatIDParam + strconv.FormatInt(chatID, 10) + telegramTextParam + message
	_, err := http.Get(sendMessageURL)
	if err != nil {
		panic(err)
	}
}
