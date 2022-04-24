package telegrambot

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"isvacbanned/util"
	"log"
)

//go:generate moq -stub -pkg mocks -out mocks/start_bot_telegram_gateway.go . StartBotTelegramGateway

type StartBotTelegramGateway interface {
	SendMessageToChat(bot *tb.Bot, chat *tb.Chat, message string)
}

type StartBotUseCase struct {
	TelegramGateway StartBotTelegramGateway
}

func NewStartBotUseCase(startBotTelegramGateway StartBotTelegramGateway) *StartBotUseCase {
	return &StartBotUseCase{TelegramGateway: startBotTelegramGateway}
}

func (uc * StartBotUseCase) SetUpStartHandler(m *tb.Message, bot *tb.Bot) {
	log.Printf("M=setUpStartHandler telegramID=%v\n", m.Chat.ID)
	uc.StartHandler(m, bot)
}

//StartHandler handles show requests
func (uc * StartBotUseCase) StartHandler(m *tb.Message, bot *tb.Bot) {
	message := util.GetStartResponse(m.Sender.FirstName)


	uc.TelegramGateway.SendMessageToChat(bot, m.Chat, message)
}
