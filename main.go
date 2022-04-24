package main

import (
	"context"
	"isvacbanned/domain/job"
	"isvacbanned/domain/telegrambot"
	"isvacbanned/gateways/mysql"
	"isvacbanned/gateways/steam"
	"isvacbanned/gateways/telegram"
	"isvacbanned/handler"
	"isvacbanned/util"
	"log"
	"net/http"
	"os"
	"strconv"

	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	startTelegramBot()
}

type TelegramBotHandler struct {
	*telegrambot.FollowPlayerUseCase
	*telegrambot.ShowPlayersUseCase
	*telegrambot.ShowSummaryUseCase
	*telegrambot.StartBotUseCase
	*telegrambot.StopNotificationsUseCase
	*telegrambot.UnfollowPlayerUseCase
}

func NewTelegramBotHandler(
	startBot *telegrambot.StartBotUseCase,
	followPlayer *telegrambot.FollowPlayerUseCase,
	unfollowPlayer *telegrambot.UnfollowPlayerUseCase,
	showPlayers *telegrambot.ShowPlayersUseCase,
	showSummary *telegrambot.ShowSummaryUseCase,
	stopNotification *telegrambot.StopNotificationsUseCase,
	) *TelegramBotHandler {
	return &TelegramBotHandler{
		FollowPlayerUseCase:      followPlayer,
		ShowPlayersUseCase:       showPlayers,
		ShowSummaryUseCase:       showSummary,
		StartBotUseCase:          startBot,
		StopNotificationsUseCase: stopNotification,
		UnfollowPlayerUseCase:    unfollowPlayer,
	}
}

func startTelegramBot() {
	var (
		port      = os.Getenv("PORT")
		publicURL = os.Getenv("PUBLIC_URL")
		token     = os.Getenv("TOKEN")
	)

	checkBanInterval, err := strconv.Atoi(os.Getenv("BAN_CHECK_INTERVAL"))
	if err != nil {
		log.Printf(`M=startTelegramBot fail to retrieve env var BAN_CHECK_INTERVAL fallback to 60 err:%s`, err.Error())
		checkBanInterval = 60
	}

	checkNicknameInterval, err := strconv.Atoi(os.Getenv("NICKNAME_CHECK_INTERVAL"))
	if err != nil {
		log.Printf(`M=startTelegramBot fail to retrieve env var NICKNAME_CHECK_INTERVAL fallback to 10 err:%s`, err.Error())
		checkNicknameInterval = 10
	}

	ctx := context.Background()

	database := util.StartDatabase()

	webhook := &tb.Webhook{
		Listen:   ":" + port,
		Endpoint: &tb.WebhookEndpoint{PublicURL: publicURL},
	}

	telegramGateway := &telegram.Messenger{}
	startBot := telegrambot.NewStartBotUseCase(telegramGateway)
	steamGateway := &steam.Steam{HTTPClient: &http.Client{}}
	userPersistenceGateway := &mysql.User{DB: database}
	followPersistenceGateway := &mysql.Follow{DB: database}
	followPlayer := telegrambot.NewFollowPlayerUseCase(steamGateway, userPersistenceGateway, followPersistenceGateway, telegramGateway)
	unfollowPlayer := telegrambot.NewUnfollowPlayerUseCase(steamGateway, userPersistenceGateway, followPersistenceGateway, telegramGateway)
	showPLayers := telegrambot.NewShowPlayersUseCase(userPersistenceGateway, followPersistenceGateway, telegramGateway)
	showSummary := telegrambot.NewShowSummaryUseCase(followPersistenceGateway, userPersistenceGateway, telegramGateway)
	stopNotifications := telegrambot.NewStopNotificationsUseCase(userPersistenceGateway, telegramGateway)

	telegramBotHandler := NewTelegramBotHandler(startBot, followPlayer, unfollowPlayer, showPLayers, showSummary, stopNotifications)

	followedUsersJobUseCase := job.NewFollowedUsersJobUseCase(followPersistenceGateway, steamGateway, telegramGateway)

	go followedUsersJobUseCase.RunScheduler(ctx, uint64(checkBanInterval), uint64(checkNicknameInterval))

	setUpBot(webhook, token, telegramBotHandler)
}

// setUpBot sets up the bot configs and its handlers
func setUpBot(webhook *tb.Webhook, token string, telegramBotHandler *TelegramBotHandler) {
	pref := tb.Settings{
		Token:  token,
		Poller: webhook,
	}

	bot, err := tb.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	setUpBotHandlers(bot, telegramBotHandler)

	bot.Start()
}

func setUpBotHandlers(bot *tb.Bot, telegramBotHandler *TelegramBotHandler) {
	bot.Handle("/follow", func(m *tb.Message) { telegramBotHandler.SetUpFollowHandler(m, bot) })

	bot.Handle("/show", func(m *tb.Message) { telegramBotHandler.SetUpShowHandler(m, bot) })

	bot.Handle("/unfollow", func(m *tb.Message) { telegramBotHandler.SetUpUnfollowHandler(m, bot) })

	bot.Handle("/stop", func(m *tb.Message) { telegramBotHandler.SetUpStopHandler(m, bot) })

	bot.Handle("/start", func(m *tb.Message) { telegramBotHandler.SetUpStartHandler(m, bot) })

	bot.Handle("/summary", func(m *tb.Message) { telegramBotHandler.SetUpSummaryHandler(m, bot) })

	bot.Handle("/CSouVALO", func(m *tb.Message) { setUpChooserHandler(m, bot) })
}

func setUpChooserHandler(m *tb.Message, bot *tb.Bot) {
	log.Printf("M=setUpChooserHandler telegramID=%v\n", m.Chat.ID)

	handler.HandleChooserRequest(m, bot)
}