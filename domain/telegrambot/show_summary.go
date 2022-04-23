package telegrambot

import (
	"context"
	"database/sql"
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"isvacbanned/util"
	"log"
)

//go:generate moq -stub -pkg mocks -out mocks/show_summary_user_persistence_gateway.go . ShowSummaryUserPersistenceGateway
//go:generate moq -stub -pkg mocks -out mocks/show_summary_follow_persistence_gateway.go . ShowSummaryFollowPersistenceGateway
//go:generate moq -stub -pkg mocks -out mocks/show_summary_telegram_gateway.go . ShowSummaryTelegramGateway

type ShowSummaryUserPersistenceGateway interface {
	GetUserID(ctx context.Context, telegramID int64) (int64, error)
	CreateUser(ctx context.Context, firstName, username string, telegramID int64) (int64, error)
}

type ShowSummaryFollowPersistenceGateway interface {
	GetUsersFollowedSummary(ctx context.Context, userID int64) (map[bool]int, error)
}

type ShowSummaryTelegramGateway interface {
	SendMessageToChat(bot *tb.Bot, chat *tb.Chat, message string)
}

type ShowSummaryUseCase struct {
	UserPersistenceGateway ShowSummaryUserPersistenceGateway
	FollowPersistenceGateway ShowSummaryFollowPersistenceGateway
	TelegramGateway ShowSummaryTelegramGateway
}

func NewShowSummaryUseCase(
	followPersistenceGateway ShowSummaryFollowPersistenceGateway,
	userPersistenceGateway ShowSummaryUserPersistenceGateway,
	telegramGateway ShowSummaryTelegramGateway) *ShowSummaryUseCase {

	return &ShowSummaryUseCase{
		UserPersistenceGateway: userPersistenceGateway,
		FollowPersistenceGateway: followPersistenceGateway,
		TelegramGateway: telegramGateway,
	}
}

func (uc *ShowSummaryUseCase) SetUpSummaryHandler(m *tb.Message, bot *tb.Bot) {
	log.Printf("M=setUpSummaryHandler telegramID=%v\n", m.Chat.ID)
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

	uc.HandleSummaryRequest(ctx, m, bot, userID)
}

//HandleSummaryRequest handles show requests
func (uc *ShowSummaryUseCase) HandleSummaryRequest(ctx context.Context, m *tb.Message, bot *tb.Bot, userID int64) {
	log.Printf("M=HandleSummaryRequest L=I userID=%d \n", m.Chat.ID)

	summary, err := uc.FollowPersistenceGateway.GetUsersFollowedSummary(ctx, userID)

	if err != nil {
		log.Printf("M=HandleSummaryRequest L=E error fetching followed users summary userID=%d err=%s \n", m.Chat.ID, err.Error())
	}

	uc.TelegramGateway.SendMessageToChat(bot, m.Chat, util.GetSummaryResponse(summary))
}
