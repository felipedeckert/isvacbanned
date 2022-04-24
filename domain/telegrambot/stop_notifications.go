package telegrambot

import (
	"context"
	"database/sql"
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"isvacbanned/util"
	"log"
)

//go:generate moq -stub -pkg mocks -out mocks/stop_notifications_persistence_gateway.go . StopNotificationsUserPersistenceGateway
//go:generate moq -stub -pkg mocks -out mocks/stop_notifications_telegram_gateway.go . StopNotificationsTelegramGateway

type StopNotificationsUserPersistenceGateway interface {
	GetUserID(ctx context.Context, telegramID int64) (int64, error)
	CreateUser(ctx context.Context, firstName, username string, telegramID int64) (int64, error)
	InactivateUser(ctx context.Context, userID int64) (int64, error)
}

type StopNotificationsTelegramGateway interface {
	SendMessageToChat(bot *tb.Bot, chat *tb.Chat, message string)
}

type StopNotificationsUseCase struct {
	PersistenceGateway StopNotificationsUserPersistenceGateway
	TelegramGateway StopNotificationsTelegramGateway
}

func NewStopNotificationsUseCase(
	persistenceGateway StopNotificationsUserPersistenceGateway,
	telegramGateway StartBotTelegramGateway) * StopNotificationsUseCase {
	return &StopNotificationsUseCase{
		PersistenceGateway: persistenceGateway,
		TelegramGateway: telegramGateway,
	}
}

func (uc * StopNotificationsUseCase) SetUpStopHandler(m *tb.Message, bot *tb.Bot) {
	log.Printf("M=setUpStopHandler telegramID=%v\n", m.Chat.ID)
	ctx := context.Background()
	userID, err := uc.PersistenceGateway.GetUserID(ctx, m.Chat.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			userID, err = uc.PersistenceGateway.CreateUser(ctx, m.Chat.FirstName, m.Chat.Username, m.Chat.ID)
			if err != nil {
				log.Printf(fmt.Errorf(`M=setUpFollowHandler error creating user err:%w`, err).Error())
				return
			}
		} else {
			log.Printf(fmt.Errorf(`M=setUpFollowHandler error getting user err:%w`, err).Error())
			return
		}
	}

	uc.StopHandler(ctx, m, bot, userID)
}

//StopHandler handles show requests
func (uc * StopNotificationsUseCase) StopHandler(ctx context.Context, m *tb.Message, bot *tb.Bot, userID int64) {
	_, err := uc.PersistenceGateway.InactivateUser(ctx, userID)
	if err != nil {
		log.Printf("M=StopHandler error inactivating user userID=%d err:%s", userID, err.Error())
		return
	}

	uc.TelegramGateway.SendMessageToChat(bot, m.Chat, util.GetStopResponse())
}