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

func TestShowSummaryUseCase_SetUpSummaryHandler_Success(t *testing.T) {
	followPersistenceGatewayMock := &mocks.ShowSummaryFollowPersistenceGatewayMock{
		GetUsersFollowedSummaryFunc: func(ctx context.Context, userID int64) (map[bool]int, error) {
			myMap := make(map[bool]int)

			return myMap, nil
		},
	}

	userPersistenceGatewayMock := &mocks.ShowSummaryUserPersistenceGatewayMock{
		CreateUserFunc: func(ctx context.Context, firstName string, username string, telegramID int64) (int64, error) {
			return 123, nil
		},
		GetUserIDFunc: func(ctx context.Context, telegramID int64) (int64, error) {
			return -1, sql.ErrNoRows
		},
	}

	actualTelegramCalls := 0
	telegramGateway := &mocks.ShowSummaryTelegramGatewayMock{SendMessageToChatFunc: func(bot *tb.Bot, chat *tb.Chat, message string) {
		actualTelegramCalls++
	}}

	useCase := telegrambot.NewShowSummaryUseCase(followPersistenceGatewayMock, userPersistenceGatewayMock, telegramGateway)

	message := &tb.Message{
		Chat: &tb.Chat{ID: 1234},
	}

	useCase.SetUpSummaryHandler(message, &tb.Bot{})
	require.Equal(t, 1, actualTelegramCalls)
}