package telegrambot

import (
	"context"
	"database/sql"
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"isvacbanned/domain/entities"
	"log"
	"strings"
)

//go:generate moq -stub -pkg mocks -out mocks/show_players_user_persistence_gateway.go . ShowPlayersUserPersistenceGateway
//go:generate moq -stub -pkg mocks -out mocks/show_players_follow_persistence_gateway.go . ShowPlayersFollowPersistenceGateway
//go:generate moq -stub -pkg mocks -out mocks/show_players_telegram_gateway.go . ShowPlayersTelegramGateway

const maxMessageLength int = 4000

type ShowPlayersUserPersistenceGateway interface {
	GetUserID(ctx context.Context, telegramID int64) (int64, error)
	CreateUser(ctx context.Context, firstName, username string, telegramID int64) (int64, error)
}

type ShowPlayersFollowPersistenceGateway interface {
	GetUsersFollowed(ctx context.Context, userID int64) ([]entities.UsersFollowed, error)
}

type ShowPlayersTelegramGateway interface {
	SendMessageToChat(bot *tb.Bot, chat *tb.Chat, message string)
}

type ShowPlayersUseCase struct {
	UserPersistenceGateway   ShowPlayersUserPersistenceGateway
	FollowPersistenceGateway ShowPlayersFollowPersistenceGateway
	TelegramGateway          ShowPlayersTelegramGateway
}

func NewShowPlayersUseCase(
	userPersistenceGateway ShowPlayersUserPersistenceGateway,
	followPersistenceGateway ShowPlayersFollowPersistenceGateway,
	telegramGateway ShowPlayersTelegramGateway) *ShowPlayersUseCase{
	return &ShowPlayersUseCase{
		UserPersistenceGateway: userPersistenceGateway,
		FollowPersistenceGateway: followPersistenceGateway,
		TelegramGateway: telegramGateway,
	}
}

func (uc *ShowPlayersUseCase) SetUpShowHandler(m *tb.Message, bot *tb.Bot) {
	log.Printf("M=setUpShowHandler telegramID=%v\n", m.Chat.ID)
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

	uc.handleShowCommand(ctx, m, bot, userID)
}

//ShowHandler handles show requests
func (uc *ShowPlayersUseCase) handleShowCommand(ctx context.Context, m *tb.Message, bot *tb.Bot, userID int64) {
	followedUsers, err := uc.FollowPersistenceGateway.GetUsersFollowed(ctx, userID)
	if err != nil {
		log.Printf("M=handleShowCommand L=E userID=%d err:%s\n", m.Chat.ID, err.Error())
		return
	}
	log.Printf("M=handleShowCommand L=I userID=%v usersFollowedCount=%v\n", m.Chat.ID, len(followedUsers))

	uc.sendShowResponse(followedUsers, bot, m)
}

func (uc *ShowPlayersUseCase) sendShowResponse(followedUsers []entities.UsersFollowed, bot *tb.Bot, m *tb.Message) {
	total := len(followedUsers)

	if total == 0 {
		uc.TelegramGateway.SendMessageToChat(bot, m.Chat, "You're not following any player yet!")
	}

	uc.sendMessageBatch(followedUsers, bot, m)
}

func (uc *ShowPlayersUseCase) sendMessageBatch(followedUsers []entities.UsersFollowed, bot *tb.Bot, m *tb.Message) {
	var sb strings.Builder
	var count int
	var status string
	var prefix string
	sb.WriteString("You're following these users: \n")
	for _, user := range followedUsers {
		count++
		prefix = ""
		status = "NOT BANNED"
		if user.IsCompleted {
			prefix = "❌ "
			status = "BANNED ❌"
		}
		sb.WriteString(prefix + user.OldNickname + " : " + status + ",\n")
		if sb.Len() > maxMessageLength {
			uc.TelegramGateway.SendMessageToChat(bot, m.Chat, sb.String())
			sb.Reset()
		}
	}

	if sb.Len() > 0 {
		uc.TelegramGateway.SendMessageToChat(bot, m.Chat, sb.String())
	}
}
