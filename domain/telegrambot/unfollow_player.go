package telegrambot

import (
	"context"
	"database/sql"
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"isvacbanned/util"
	"log"
)

//go:generate moq -stub -pkg mocks -out mocks/unfollow_player_steam_gateway.go . UnfollowPlayerSteamGateway
//go:generate moq -stub -pkg mocks -out mocks/unfollow_player_user_persistence_gateway.go . UnfollowPlayerUserPersistenceGateway
//go:generate moq -stub -pkg mocks -out mocks/unfollow_player_follow_persistence_gateway.go . UnfollowPlayerFollowPersistenceGateway
//go:generate moq -stub -pkg mocks -out mocks/unfollow_player_telegram_gateway.go . UnfollowPlayerTelegramGateway

type UnfollowPlayerSteamGateway interface {
	GetPlayerSteamID(playerName string) (string, error)
}

type UnfollowPlayerUserPersistenceGateway interface {
	GetUserID(ctx context.Context, telegramID int64) (int64, error)
	CreateUser(ctx context.Context, firstName, username string, telegramID int64) (int64, error)
}

type UnfollowPlayerFollowPersistenceGateway interface {
	UnfollowSteamUser(ctx context.Context, userID int64, steamID string) (int64, error)
}

type UnfollowPlayerTelegramGateway interface {
	SendMessageToChat(bot *tb.Bot, chat *tb.Chat, message string)
}

type UnfollowPlayerUseCase struct {
	SteamGateway             UnfollowPlayerSteamGateway
	UserPersistenceGateway   UnfollowPlayerUserPersistenceGateway
	FollowPersistenceGateway UnfollowPlayerFollowPersistenceGateway
	TelegramGateway			 UnfollowPlayerTelegramGateway
}

func NewUnfollowPlayerUseCase (
	steamGateway UnfollowPlayerSteamGateway,
	userPersistenceGateway UnfollowPlayerUserPersistenceGateway,
	followPersistenceGateway UnfollowPlayerFollowPersistenceGateway,
	telegramGateway			 UnfollowPlayerTelegramGateway) *UnfollowPlayerUseCase{
	return &UnfollowPlayerUseCase{
		SteamGateway:           steamGateway,
		UserPersistenceGateway: userPersistenceGateway,
		FollowPersistenceGateway: followPersistenceGateway,
		TelegramGateway: telegramGateway,
	}
}

func (uc * UnfollowPlayerUseCase) SetUpUnfollowHandler(m *tb.Message, bot *tb.Bot) {
	ctx := context.Background()
	userID, err := uc.UserPersistenceGateway.GetUserID(ctx, m.Chat.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			userID, err = uc.UserPersistenceGateway.CreateUser(ctx, m.Chat.FirstName, m.Chat.Username, m.Chat.ID)
			if err != nil {
				log.Printf(fmt.Errorf(`M=setUpFollowHandler error creating user err:%w`, err).Error())
				return
			}
		} else {
			log.Printf(fmt.Errorf(`M=setUpFollowHandler error getting user err:%w`, err).Error())
			return
		}
	}

	steamID, err := uc.SteamGateway.GetPlayerSteamID(m.Payload)

	log.Printf("M=setUpUnfollowHandler chatID=%v steamID=%v\n", m.Chat.ID, steamID)

	if err != nil {
		log.Printf(`M=setUpUnfollowHandler err=%s`, err.Error())
		return
	} else if len(steamID) != steamIDLength || !isNumeric(steamID) {
		_, err = bot.Send(m.Chat, "Invalid Param!")
		if err != nil {
			log.Printf(`M=setUpUnfollowHandler step=send err=%s`, err.Error())
		}
		return
	}

	uc.handleUnfollowRequest(ctx, m, bot, steamID, userID)
}

//handleUnfollowRequest handles a follow request
func (uc * UnfollowPlayerUseCase) handleUnfollowRequest(ctx context.Context, m *tb.Message, bot *tb.Bot, steamID string, userID int64) {
	rows, err := uc.FollowPersistenceGateway.UnfollowSteamUser(ctx, userID, steamID)

	if err != nil {
		log.Printf(`M=handleUnfollowRequest err unfollowing user steamID=%s err:%s`, steamID, err.Error())
		return
	}

	var message string
	if rows != 1 {
		message = "Unable to unfollow this player. Are you sure you follow him/her?"
	} else {
		message = "You will NOT receive more updates about this player: " + util.SteamProfileURL + steamID
	}

	uc.TelegramGateway.SendMessageToChat(bot, m.Chat, message)
}
