// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"isvacbanned/domain/telegrambot"
	"sync"
)

// Ensure, that ShowSummaryTelegramGatewayMock does implement telegrambot.ShowSummaryTelegramGateway.
// If this is not the case, regenerate this file with moq.
var _ telegrambot.ShowSummaryTelegramGateway = &ShowSummaryTelegramGatewayMock{}

// ShowSummaryTelegramGatewayMock is a mock implementation of telegrambot.ShowSummaryTelegramGateway.
//
// 	func TestSomethingThatUsesShowSummaryTelegramGateway(t *testing.T) {
//
// 		// make and configure a mocked telegrambot.ShowSummaryTelegramGateway
// 		mockedShowSummaryTelegramGateway := &ShowSummaryTelegramGatewayMock{
// 			SendMessageToChatFunc: func(bot *tb.Bot, chat *tb.Chat, message string)  {
// 				panic("mock out the SendMessageToChat method")
// 			},
// 		}
//
// 		// use mockedShowSummaryTelegramGateway in code that requires telegrambot.ShowSummaryTelegramGateway
// 		// and then make assertions.
//
// 	}
type ShowSummaryTelegramGatewayMock struct {
	// SendMessageToChatFunc mocks the SendMessageToChat method.
	SendMessageToChatFunc func(bot *tb.Bot, chat *tb.Chat, message string)

	// calls tracks calls to the methods.
	calls struct {
		// SendMessageToChat holds details about calls to the SendMessageToChat method.
		SendMessageToChat []struct {
			// Bot is the bot argument value.
			Bot *tb.Bot
			// Chat is the chat argument value.
			Chat *tb.Chat
			// Message is the message argument value.
			Message string
		}
	}
	lockSendMessageToChat sync.RWMutex
}

// SendMessageToChat calls SendMessageToChatFunc.
func (mock *ShowSummaryTelegramGatewayMock) SendMessageToChat(bot *tb.Bot, chat *tb.Chat, message string) {
	callInfo := struct {
		Bot     *tb.Bot
		Chat    *tb.Chat
		Message string
	}{
		Bot:     bot,
		Chat:    chat,
		Message: message,
	}
	mock.lockSendMessageToChat.Lock()
	mock.calls.SendMessageToChat = append(mock.calls.SendMessageToChat, callInfo)
	mock.lockSendMessageToChat.Unlock()
	if mock.SendMessageToChatFunc == nil {
		return
	}
	mock.SendMessageToChatFunc(bot, chat, message)
}

// SendMessageToChatCalls gets all the calls that were made to SendMessageToChat.
// Check the length with:
//     len(mockedShowSummaryTelegramGateway.SendMessageToChatCalls())
func (mock *ShowSummaryTelegramGatewayMock) SendMessageToChatCalls() []struct {
	Bot     *tb.Bot
	Chat    *tb.Chat
	Message string
} {
	var calls []struct {
		Bot     *tb.Bot
		Chat    *tb.Chat
		Message string
	}
	mock.lockSendMessageToChat.RLock()
	calls = mock.calls.SendMessageToChat
	mock.lockSendMessageToChat.RUnlock()
	return calls
}