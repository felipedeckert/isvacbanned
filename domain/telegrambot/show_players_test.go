package telegrambot_test

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	tb "gopkg.in/tucnak/telebot.v2"
	"isvacbanned/domain/entities"
	"isvacbanned/domain/telegrambot"
	"isvacbanned/domain/telegrambot/mocks"
	"testing"
)

func TestShowPlayersUseCase_SetUpShowHandler_Success(t *testing.T) {
	userPersistenceGateway := &mocks.ShowPlayersUserPersistenceGatewayMock{
		CreateUserFunc: func(ctx context.Context, firstName string, username string, telegramID int64) (int64, error) {
			return 123, nil
		},
		GetUserIDFunc: func(ctx context.Context, telegramID int64) (int64, error) {
			return -1, sql.ErrNoRows
		},
	}

	followPersistenceGateway := &mocks.ShowPlayersFollowPersistenceGatewayMock{
		GetUsersFollowedFunc: func(ctx context.Context, userID int64) ([]entities.UsersFollowed, error) {
			return []entities.UsersFollowed{
				{
					ID:           10,
					SteamID:      "76561197960690195",
					OldNickname:  "fallen",
					CurrNickname: "fln",
					IsCompleted:  false,
				},
				{
					ID:           11,
					SteamID:      "76561197960690196",
					OldNickname:  "fer",
					CurrNickname: "fergod",
					IsCompleted:  false,
				},
			}, nil
		},
	}

	actualTelegramCalls := 0
	telegramGateway := &mocks.ShowPlayersTelegramGatewayMock{SendMessageToChatFunc: func(bot *tb.Bot, chat *tb.Chat, message string) {
		actualTelegramCalls++
	}}

	useCase := telegrambot.NewShowPlayersUseCase(userPersistenceGateway, followPersistenceGateway, telegramGateway)

	message := &tb.Message{
		Chat: &tb.Chat{ID: 1234},
	}

	useCase.SetUpShowHandler(message, &tb.Bot{})

	require.Equal(t, 1, actualTelegramCalls)
}