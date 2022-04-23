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

func TestFollowPlayerUseCase_SetUpFollowHandler_Success(t *testing.T) {
	steamGatewayMock := &mocks.FollowPlayerSteamGatewayMock{
		GetPlayerSteamIDFunc: func(playerName string) (string, error) {
			return "76561197960690195", nil
		},
		GetPlayersCurrentNicknamesFunc: func(steamIDs ...string) (entities.PlayerInfo, error) {
			return entities.PlayerInfo{
				Response: entities.PlayerNicknameData{
					Players: []entities.ResponseNicknameData{{
						PersonaName: "fallen",
						SteamID:     "76561197960690195",
					}},
				},
			}, nil
		},
		GetPlayersStatusFunc: func(steamIDs ...string) (entities.Player, error) {
			return entities.Player{
				Players: []entities.PlayerData{
					{
						SteamId: "76561197960690195",
						CommunityBanned: false,
						VACBanned: false,
						NumberOfVACBans: 0,
						DaysSinceLastBan: 0,
						NumberOfGameBans: 0,
						EconomyBan: "no",
					},
				},
			}, nil
		},

	}

	userPersistenceGateway := &mocks.FollowPlayerUserPersistenceGatewayMock{
		CreateUserFunc: func(ctx context.Context, firstName string, username string, telegramID int64) (int64, error) {
			return 123, nil
		},
		GetUserIDFunc: func(ctx context.Context, telegramID int64) (int64, error) {
			return 0, sql.ErrNoRows
		},
		ActivateUserFunc: func(ctx context.Context, userID int64) (int64, error) {
			return 123, nil
		},
	}

	actualSetFollowedUserToCompletedCalls := 0
	followPersistenceGatewayMock := &mocks.FollowPlayerFollowPersistenceGatewayMock{
		FollowSteamUserFunc: func(ctx context.Context, chatID int64, steamID string, currNickname string, userID int64) (int64, error) {
			return 321, nil
		},
		GetFollowerCountBySteamIDFunc: func(ctx context.Context, steamID string) (int64, error) {
			return 2, nil
		},
		GetUsersFollowedFunc:           nil,
		IsFollowedFunc: func(ctx context.Context, steamID string, userID int64) (string, int64, error) {
			return "", -1, nil
		},
		SetFollowedUserToCompletedFunc: func(ctx context.Context, id []int64) {
			actualSetFollowedUserToCompletedCalls++
		},
	}
	actualTelegramCalls := 0
	telegramGateway := &mocks.FollowPlayerTelegramGatewayMock{
		SendMessageToChatFunc: func(bot *tb.Bot, chat *tb.Chat, message string) {
		actualTelegramCalls++
	}}

	useCase := telegrambot.NewFollowPlayerUseCase(steamGatewayMock, userPersistenceGateway, followPersistenceGatewayMock, telegramGateway)

	message := &tb.Message{
		Chat: &tb.Chat{ID: 1234},
	}
	useCase.SetUpFollowHandler(message, &tb.Bot{})

	require.Equal(t, 1, actualTelegramCalls)
	require.Equal(t, 0, actualSetFollowedUserToCompletedCalls)
}
