package messenger

import tb "gopkg.in/tucnak/telebot.v2"

//MessengerMock is the mock implementation of the Messenger
type MessengerMock struct{}

//SendMessage is the mock implementation of SendMessage
func (m MessengerMock) SendMessage(bot *tb.Bot, user *tb.User, message string) {}

//SendMessageToChat is the mock implementation of SendMessageToChat
func (m MessengerMock) SendMessageToChat(bot *tb.Bot, chat *tb.Chat, message string) {}

//SendMessageToUser is the mock implementation of SendMessageToUser
func (m MessengerMock) SendMessageToUser(message string, chatID int64) {}
