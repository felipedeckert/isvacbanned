package telegrambot_test

import (
	"github.com/stretchr/testify/require"
	tb "gopkg.in/tucnak/telebot.v2"
	"isvacbanned/domain/telegrambot"
	"isvacbanned/domain/telegrambot/mocks"
	"testing"
)

func TestStartBotUseCase_SetUpStartHandler_Success(t *testing.T) {
	actualTelegramCalls := 0
	telegramGateway := &mocks.StartBotTelegramGatewayMock{SendMessageToChatFunc: func(bot *tb.Bot, chat *tb.Chat, message string) {
		actualTelegramCalls++
	}}

	useCase := telegrambot.NewStartBotUseCase(telegramGateway)

	message := &tb.Message{
		Chat: &tb.Chat{ID: 1234},
		Sender: &tb.User{FirstName: "firstName"},
	}

	useCase.SetUpStartHandler(message, &tb.Bot{})
	require.Equal(t, 1, actualTelegramCalls)
}