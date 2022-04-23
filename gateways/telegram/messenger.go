package telegram

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const telegramAPIURL = "https://api.telegram.org/bot"
const telegramMethod = "/sendMessage"
const telegramChatIDParam = "?chat_id="
const telegramTextParam = "&text="

//MessengerInterface is the interface that defines the behavior of the Messenger
type MessengerInterface interface {
	SendMessage(bot *tb.Bot, user *tb.User, message string)
	SendMessageToChat(bot *tb.Bot, chat *tb.Chat, message string)
	SendMessageToUser(message string, chatID int64)
}

//Messenger is the real implementation of the Messenger
type Messenger struct{}

//SendMessage sends a message to user via bot
func (m *Messenger) SendMessage(bot *tb.Bot, user *tb.User, message string) {
	log.Println("M=SendMessage step=start")
	_, err := bot.Send(user, message)
	if err != nil {
		log.Printf("M=SendMessage err=%s", err.Error())
	}
	log.Println("M=SendMessage step=end")
}

//SendMessageToChat sends a message to chat via bot
func (m *Messenger) SendMessageToChat(bot *tb.Bot, chat *tb.Chat, message string) {
	log.Println("M=SendMessageToChat step=start")
	_, err := bot.Send(chat, message)
	if err != nil {
		log.Printf("M=SendMessageToChat err=%s", err.Error())
	}
	log.Println("M=SendMessageToChat step=end")
}

//SendMessageToUser sends a message to user/group via telegram API
func (m *Messenger) SendMessageToUser(message string, chatID int64) {
	log.Println("M=SendMessageToUser step=start")
	token := os.Getenv("TOKEN")

	sendMessageURL := telegramAPIURL + token + telegramMethod + telegramChatIDParam + strconv.FormatInt(chatID, 10) + telegramTextParam + url.QueryEscape(message)

	_, err := http.Get(sendMessageURL)
	if err != nil {
		log.Printf(`M=SendMessageToUser L=E error sending message URL=%s, err=%s`, sendMessageURL, err.Error())
	}
	log.Println("M=SendMessageToUser step=end")
}