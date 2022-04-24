package telegrambot

import (
	"context"
	"database/sql"
	tb "gopkg.in/tucnak/telebot.v2"
	"isvacbanned/domain/entities"
	"isvacbanned/util"
	"log"
	"strconv"
)

//go:generate moq -stub -pkg mocks -out mocks/follow_player_steam_gateway.go . FollowPlayerSteamGateway
//go:generate moq -stub -pkg mocks -out mocks/follow_player_user_persistence_gateway.go . FollowPlayerUserPersistenceGateway
//go:generate moq -stub -pkg mocks -out mocks/follow_player_follow_persistence_gateway.go . FollowPlayerFollowPersistenceGateway
//go:generate moq -stub -pkg mocks -out mocks/follow_player_telegram_gateway.go . FollowPlayerTelegramGateway

const steamIDLength = 17

type FollowPlayerSteamGateway interface {
	GetSteamID(param string) (string, error)
	GetPlayerSteamID(playerName string) (string, error)
	GetPlayersCurrentNicknames(steamIDs ...string) (entities.PlayerInfo, error)
	GetPlayersStatus(steamIDs ...string) (entities.Player, error)
}

type FollowPlayerUserPersistenceGateway interface {
	GetUserID(ctx context.Context, telegramID int64) (int64, error)
	CreateUser(ctx context.Context, firstName, username string, telegramID int64) (int64, error)
	ActivateUser(ctx context.Context, userID int64) (int64, error)
}

type FollowPlayerFollowPersistenceGateway interface {
	GetFollowerCountBySteamID(ctx context.Context, steamID string) (int64, error)
	GetUsersFollowed(ctx context.Context, userID int64) ([]entities.UsersFollowed, error)
	IsFollowed(ctx context.Context, steamID string, userID int64) (string, int64, error)
	FollowSteamUser(ctx context.Context, chatID int64, steamID, currNickname string, userID int64) (int64, error)
	SetFollowedUserToCompleted(ctx context.Context, id []int64)
}

type FollowPlayerTelegramGateway interface {
	SendMessageToChat(bot *tb.Bot, chat *tb.Chat, message string)
}

type FollowPlayerUseCase struct {
	SteamGateway             FollowPlayerSteamGateway
	UserPersistenceGateway   FollowPlayerUserPersistenceGateway
	FollowPersistenceGateway FollowPlayerFollowPersistenceGateway
	TelegramGateway			 FollowPlayerTelegramGateway
}

func NewFollowPlayerUseCase(
	steamGateway 			 FollowPlayerSteamGateway,
	userPersistenceGateway   FollowPlayerUserPersistenceGateway,
	followPersistenceGateway FollowPlayerFollowPersistenceGateway,
	telegramGateway			 FollowPlayerTelegramGateway) *FollowPlayerUseCase{
	return &FollowPlayerUseCase{
		SteamGateway:           steamGateway,
		UserPersistenceGateway: userPersistenceGateway,
		FollowPersistenceGateway: followPersistenceGateway,
		TelegramGateway: telegramGateway,
	}
}

func (uc * FollowPlayerUseCase) SetUpFollowHandler(m *tb.Message, bot *tb.Bot) int64 {
	ctx := context.Background()
	userID, err := uc.UserPersistenceGateway.GetUserID(ctx, m.Chat.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			userID, err = uc.UserPersistenceGateway.CreateUser(ctx, m.Chat.FirstName, m.Chat.Username, m.Chat.ID)
			if err != nil {
				log.Printf(`M=setUpFollowHandler error creating user err:%s`, err.Error())
				return -1
			}
		} else {
			log.Printf(`M=setUpFollowHandler error getting user err:%s`, err.Error())
			return -1
		}
	}

	steamID, err := uc.SteamGateway.GetSteamID(m.Payload)
	log.Printf("M=setUpFollowHandler payload=%s userID=%d steamID=%s\n", m.Payload, userID, steamID)

	if err != nil {
		log.Printf(`M=setUpFollowHandler err=%s`, err.Error())
		return -1
	} else if len(steamID) != steamIDLength || !isNumeric(steamID) {
		_, err = bot.Send(m.Chat, "Invalid Param!")
		if err != nil {
			log.Printf(`M=setUpFollowHandler step=send err=%s`, err.Error())
			return -1
		}
		return -1
	}

	playerInfo, err := uc.SteamGateway.GetPlayersCurrentNicknames(steamID)
	if err != nil {
		log.Printf(`M=setUpFollowHandler error while getting player's nickname err:%s`, err.Error())
		return -1
	}
	player, err := uc.SteamGateway.GetPlayersStatus(steamID)
	if err != nil {
		log.Printf(`M=setUpFollowHandler error while getting player's status err:%s`, err.Error())
		return -1
	}

	currNickname := playerInfo.Response.Players[0].PersonaName

	if len(player.Players) > 0 && currNickname != "" {
		playerData := player.Players[0]

		id, err :=  uc.HandleFollowRequest(ctx, m, bot, steamID, currNickname, userID, playerData.VACBanned)
		if err != nil {
			log.Printf(`M=setUpFollowHandler error while following player err:%s`, err.Error())
			return -1
		}

		return id
	}
	return -1
}

//HandleFollowRequest handles a follow request
func (uc *FollowPlayerUseCase) HandleFollowRequest(ctx context.Context, m *tb.Message, bot *tb.Bot, steamID, currNickname string, userID int64, isVACBanned bool) (int64, error) {
	log.Printf("M=HandleFollowRequest steamID=%s\n", steamID)
	followersCount, err := uc.FollowPersistenceGateway.GetFollowerCountBySteamID(ctx, steamID)

	_, err = uc.UserPersistenceGateway.ActivateUser(ctx, userID)
	if err != nil {
		return -1, err
	}

	var dbID int64
	oldNickname, followID, err := uc.FollowPersistenceGateway.IsFollowed(ctx, steamID, userID)
	if err != nil && err != sql.ErrNoRows {
		return -1, err
	}

	updateID := followID

	if oldNickname == "" {
		dbID, err = uc.FollowPersistenceGateway.FollowSteamUser(ctx, m.Chat.ID, steamID, currNickname, userID)
		updateID = dbID
	}

	if isVACBanned {
		uc.FollowPersistenceGateway.SetFollowedUserToCompleted(ctx, []int64{updateID})
	}

	response := util.GetFollowResponseMessage(oldNickname, currNickname, followersCount, isVACBanned)

	uc.TelegramGateway.SendMessageToChat(bot, m.Chat, response)

	return dbID, nil
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
