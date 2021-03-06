package service

import (
	"isvacbanned/handler"
	"isvacbanned/util"
	"log"
	"strconv"

	tb "gopkg.in/tucnak/telebot.v2"
)

// SetUpBot sets up the bot configs and its handlers
func SetUpBot(webhook *tb.Webhook, token string) {
	pref := tb.Settings{
		Token:  token,
		Poller: webhook,
	}

	bot, err := tb.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	setUpBotHandlers(bot)

	bot.Start()
}

func setUpBotHandlers(bot *tb.Bot) {
	bot.Handle("/follow", func(m *tb.Message) { setUpFollowHandler(m, bot) })

	bot.Handle("/show", func(m *tb.Message) { setUpShowHandler(m, bot) })

	bot.Handle("/unfollow", func(m *tb.Message) { setUpUnfollowHandler(m, bot) })

	bot.Handle("/stop", func(m *tb.Message) { setUpStopHandler(m, bot) })

	bot.Handle("/start", func(m *tb.Message) { setUpStartHandler(m, bot) })

	bot.Handle("/summary", func(m *tb.Message) { setUpSummaryHandler(m, bot) })

	bot.Handle("/CSouVALO", func(m *tb.Message) { setUpChooserHandler(m, bot) })
}

func setUpFollowHandler(m *tb.Message, bot *tb.Bot) int64 {
	userID := UserServiceClient.getUserID(m.Chat)
	steamID, err := UrlServiceClient.getSteamID(m.Payload)

	log.Printf("M=setUpFollowHandler payload=%v userID=%v steamID=%v\n", m.Payload, userID, steamID)

	if err != nil {
		log.Printf(`M=setUpFollowHandler err=%s`, err.Error())
		return -1
	} else if len(steamID) != util.SteamIDLength || !isNumeric(steamID) {
		_, err = bot.Send(m.Chat, "Invalid Param!")
		if err != nil {
			log.Printf(`M=setUpFollowHandler step=send err=%s`, err.Error())
			return 0
		}
		return -1
	}

	currNickname := PlayerServiceClient.GetPlayerCurrentNickname(steamID)
	player := PlayerServiceClient.GetPlayerStatus(steamID)

	if len(player.Players) > 0 && currNickname != "" {
		playerData := player.Players[0]

		return handler.FollowHandlerClient.HandleFollowRequest(m, bot, steamID, currNickname, userID, playerData.VACBanned)
	}
	return -1
}

func setUpShowHandler(m *tb.Message, bot *tb.Bot) {
	log.Printf("M=setUpShowHandler telegramID=%v\n", m.Chat.ID)
	userID := UserServiceClient.getUserID(m.Chat)

	handler.ShowHandler(m, bot, userID)
}

func setUpStopHandler(m *tb.Message, bot *tb.Bot) {
	log.Printf("M=setUpStopHandler telegramID=%v\n", m.Chat.ID)
	userID := UserServiceClient.getUserID(m.Chat)

	handler.StopHandler(m, bot, userID)
}

func setUpStartHandler(m *tb.Message, bot *tb.Bot) {
	log.Printf("M=setUpStartHandler telegramID=%v\n", m.Chat.ID)
	handler.StartHandler(m, bot)
}

func setUpUnfollowHandler(m *tb.Message, bot *tb.Bot) {
	userID := UserServiceClient.getUserID(m.Chat)
	steamID, err := UrlServiceClient.getSteamID(m.Payload)

	log.Printf("M=setUpUnfollowHandler chatID=%v steamID=%v\n", m.Chat.ID, steamID)

	if err != nil {
		log.Printf(`M=setUpUnfollowHandler err=%s`, err.Error())
		return
	} else if len(steamID) != util.SteamIDLength || !isNumeric(steamID) {
		_, err = bot.Send(m.Chat, "Invalid Param!")
		if err != nil {
			log.Printf(`M=setUpUnfollowHandler step=send err=%s`, err.Error())
		}
		return
	}

	handler.UnfollowHandlerClient.HandleUnfollowRequest(m, bot, steamID, userID)
}

func setUpSummaryHandler(m *tb.Message, bot *tb.Bot) {
	log.Printf("M=setUpSummaryHandler telegramID=%v\n", m.Chat.ID)
	userID := UserServiceClient.getUserID(m.Chat)

	handler.HandleSummaryRequest(m, bot, userID)
}

func setUpChooserHandler(m *tb.Message, bot *tb.Bot) {
	log.Printf("M=setUpChooserHandler telegramID=%v\n", m.Chat.ID)

	handler.HandleChooserRequest(m, bot)
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
