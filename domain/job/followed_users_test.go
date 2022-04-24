package job_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"isvacbanned/domain/entities"
	"isvacbanned/domain/job"
	"isvacbanned/domain/job/mocks"
	"math"
	"testing"
)

func TestCheckFollowedUsersBan(t *testing.T){
	var usersTotal int64 = 150
	actualSetCurrNicknameFuncCalls := 0
	actualSetFollowedUserToCompletedFuncCalls := 0
	followPersistenceGatewayMock := &mocks.FollowedUsersJobFollowPersistenceGatewayMock{
		GetAllIncompleteFollowedUsersFunc: func(ctx context.Context) (map[int64][]entities.UsersFollowed, error) {
			return createGetAllIncompleteFollowedUsersResponse(usersTotal), nil
		},
		SetCurrNicknameFunc: func(ctx context.Context, userID int64, sanitizedActualNickname string) error {
			actualSetCurrNicknameFuncCalls++
			return nil
		},
		SetFollowedUserToCompletedFunc: func(ctx context.Context, id []int64) {
			actualSetFollowedUserToCompletedFuncCalls++
		},
	}

	banRatio := 20
	nicknameChangedRatio := 10
	actualGetPlayersStatusFuncCalls := 0
	actualGetPlayersCurrentNicknamesFuncCalls := 0
	steamGatewayMock := &mocks.FollowedUsersJobSteamGatewayMock{
		GetPlayersCurrentNicknamesFunc: func(steamIDs ...string) (entities.PlayerInfo, error) {
			actualGetPlayersCurrentNicknamesFuncCalls++
			return createGetPlayersCurrentNicknamesFuncResponse(nicknameChangedRatio, steamIDs...), nil
		},
		GetPlayersStatusFunc: func(steamIDs ...string) (entities.Player, error) {
			actualGetPlayersStatusFuncCalls++
			return createGetPlayersStatusFuncResponse(banRatio, steamIDs...), nil
		},
	}

	telegramActualCalls := 0
	telegramGatewayMock := &mocks.FollowedUsersJobTelegramGatewayMock{
		SendMessageToUserFunc: func(message string, chatID int64) {
			telegramActualCalls++
		},
	}

	useCase := job.NewFollowedUsersJobUseCase(followPersistenceGatewayMock, steamGatewayMock, telegramGatewayMock)

	useCase.CheckFollowedUsersBan(context.Background())
	getPlayersCurrentNicknamesFuncExpectedCalls := int(math.Ceil(float64(usersTotal)/100.00))
	require.Equal(t, getPlayersCurrentNicknamesFuncExpectedCalls, actualGetPlayersCurrentNicknamesFuncCalls)

	getPlayersStatusFuncExpectedCalls := int(math.Ceil(float64(usersTotal)/100.00))
	require.Equal(t, getPlayersStatusFuncExpectedCalls, actualGetPlayersStatusFuncCalls)

	expectedTelegramCalls := int(math.Floor(float64(usersTotal)/float64(banRatio))) + 1
	require.Equal(t, expectedTelegramCalls, telegramActualCalls)
}

func TestCheckFollowedUsersNickname(t *testing.T){
	var usersTotal int64 = 150
	actualSetCurrNicknameFuncCalls := 0
	actualSetFollowedUserToCompletedFuncCalls := 0
	followPersistenceGatewayMock := &mocks.FollowedUsersJobFollowPersistenceGatewayMock{
		GetAllIncompleteFollowedUsersFunc: func(ctx context.Context) (map[int64][]entities.UsersFollowed, error) {
			return createGetAllIncompleteFollowedUsersResponse(usersTotal), nil
		},
		SetCurrNicknameFunc: func(ctx context.Context, userID int64, sanitizedActualNickname string) error {
			actualSetCurrNicknameFuncCalls++
			return nil
		},
		SetFollowedUserToCompletedFunc: func(ctx context.Context, id []int64) {
			actualSetFollowedUserToCompletedFuncCalls++
		},
	}

	banRatio := 20
	nicknameChangedRatio := 10
	actualGetPlayersStatusFuncCalls := 0
	actualGetPlayersCurrentNicknamesFuncCalls := 0
	steamGatewayMock := &mocks.FollowedUsersJobSteamGatewayMock{
		GetPlayersCurrentNicknamesFunc: func(steamIDs ...string) (entities.PlayerInfo, error) {
			actualGetPlayersCurrentNicknamesFuncCalls++
			return createGetPlayersCurrentNicknamesFuncResponse(nicknameChangedRatio, steamIDs...), nil
		},
		GetPlayersStatusFunc: func(steamIDs ...string) (entities.Player, error) {
			actualGetPlayersStatusFuncCalls++
			return createGetPlayersStatusFuncResponse(banRatio, steamIDs...), nil
		},
	}

	telegramActualCalls := 0
	telegramGatewayMock := &mocks.FollowedUsersJobTelegramGatewayMock{
		SendMessageToUserFunc: func(message string, chatID int64) {
			telegramActualCalls++
		},
	}

	useCase := job.NewFollowedUsersJobUseCase(followPersistenceGatewayMock, steamGatewayMock, telegramGatewayMock)

	useCase.CheckFollowedUsersNickname(context.Background())
	getPlayersCurrentNicknamesFuncExpectedCalls := int(math.Ceil(float64(usersTotal)/100.00))
	require.Equal(t, getPlayersCurrentNicknamesFuncExpectedCalls, actualGetPlayersCurrentNicknamesFuncCalls)

	getPlayersStatusFuncExpectedCalls := 0
	require.Equal(t, getPlayersStatusFuncExpectedCalls, actualGetPlayersStatusFuncCalls)

	expectedTelegramCalls := int(math.Floor(float64(usersTotal)/float64(nicknameChangedRatio)))
	require.Equal(t, expectedTelegramCalls, telegramActualCalls)
}

func createGetPlayersStatusFuncResponse(banRatio int, steamIDs ...string) entities.Player  {
	players := make([]entities.PlayerData, len(steamIDs))

	for i, steamID := range steamIDs{
		isVACBanned := false
		if i % banRatio == 0 {
			isVACBanned = true
		}
		players[i] = entities.PlayerData{
			SteamId:          steamID,
			CommunityBanned:  false,
			VACBanned:        isVACBanned,
			NumberOfVACBans:  0,
			DaysSinceLastBan: 0,
			NumberOfGameBans: 0,
			EconomyBan:       "no",
		}
	}

	return entities.Player{Players: players}
}

func createGetPlayersCurrentNicknamesFuncResponse(nicknameChangedRatio int, steamIDs ...string) entities.PlayerInfo {
	players := make([]entities.ResponseNicknameData, len(steamIDs))
	for i, steamID := range steamIDs {
		personaName := "new_nickname"
		if (i+1) % nicknameChangedRatio == 0 {
			personaName = "diff_nickname"
		}
		players[i] = entities.ResponseNicknameData{
			PersonaName: personaName,
			SteamID:     steamID,
		}
	}

	return entities.PlayerInfo{Response: entities.PlayerNicknameData{Players: players}}
}

func createGetAllIncompleteFollowedUsersResponse(size int64) map[int64][]entities.UsersFollowed {
	users := make([]entities.UsersFollowed, size)

	for i := int64(0); i < size; i++{
		index := fmt.Sprintf("%03d", i)
		users[i] = entities.UsersFollowed{
			ID:           i,
			SteamID:      "12345678912345"+index,
			OldNickname:  "old_nickname",
			CurrNickname: "new_nickname",
			IsCompleted:  false,
		}
	}

	myMap := make(map[int64][]entities.UsersFollowed)
	chatID := int64(987)
	myMap[chatID] = users

	return myMap
}
