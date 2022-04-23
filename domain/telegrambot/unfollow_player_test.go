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

func TestUnfollowPlayerUseCase_SetUpUnfollowHandler_Success(t *testing.T) {

	steamGatewayMock := &mocks.UnfollowPlayerSteamGatewayMock{GetPlayerSteamIDFunc: func(playerName string) (string, error) {
		return "76561197960690195", nil
	}}

	userPersistenceGateway := &mocks.UnfollowPlayerUserPersistenceGatewayMock{
		CreateUserFunc: func(ctx context.Context, firstName string, username string, telegramID int64) (int64, error) {
			return 123, nil
		},
		GetUserIDFunc: func(ctx context.Context, telegramID int64) (int64, error) {
			return 0, sql.ErrNoRows
		},
	}

	followPersistenceGateway := &mocks.UnfollowPlayerFollowPersistenceGatewayMock{UnfollowSteamUserFunc:
		func(ctx context.Context, userID int64, steamID string) (int64, error) {
			return 1234, nil
		},
	}

	actualTelegramCalls := 0
	telegramGateway := &mocks.UnfollowPlayerTelegramGatewayMock{SendMessageToChatFunc: func(bot *tb.Bot, chat *tb.Chat, message string) {
		actualTelegramCalls++
	}}

	useCase := telegrambot.NewUnfollowPlayerUseCase(steamGatewayMock, userPersistenceGateway, followPersistenceGateway, telegramGateway)

	message := &tb.Message{
		Chat: &tb.Chat{ID: 1234},
	}
	useCase.SetUpUnfollowHandler(message, &tb.Bot{})
	require.Equal(t, 1, actualTelegramCalls)
}