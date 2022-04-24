package telegrambot_test

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	tb "gopkg.in/tucnak/telebot.v2"
	"isvacbanned/domain/telegrambot"
	"isvacbanned/domain/telegrambot/mocks"
	"testing"
)

func TestStopNotificationsUseCase_SetUpStopHandler_Success(t *testing.T) {

	userPersistenceGatewayMock := &mocks.StopNotificationsUserPersistenceGatewayMock{
		CreateUserFunc: func(ctx context.Context, firstName string, username string, telegramID int64) (int64, error) {
			return 123, nil
		},
		GetUserIDFunc: func(ctx context.Context, telegramID int64) (int64, error) {
			return -1, sql.ErrNoRows
		},
		InactivateUserFunc: func(ctx context.Context, userID int64) (int64, error) {
			return 123, nil
		},
	}

	actualTelegramCalls := 0
	telegramGatewayMock := &mocks.StopNotificationsTelegramGatewayMock{SendMessageToChatFunc: func(bot *tb.Bot, chat *tb.Chat, message string) {
		actualTelegramCalls++
	}}

	useCase := telegrambot.NewStopNotificationsUseCase(userPersistenceGatewayMock, telegramGatewayMock)

	message := &tb.Message{
		Chat: &tb.Chat{ID: 1234},
		Sender: &tb.User{FirstName: "firstName"},
	}

	useCase.SetUpStopHandler(message, &tb.Bot{})
	require.Equal(t, 1, actualTelegramCalls)
}
