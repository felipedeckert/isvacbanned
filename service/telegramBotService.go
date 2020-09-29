package service

import (
	"isvacbanned/handler"
	"log"

	tb "gopkg.in/tucnak/telebot.v2"
)

const steamIDLength = 17
const playerBanned = "This player is VAC banned!"
const playerNotBanned = "This player is NOT VAC banned!"

var followHandler *handler.FollowHandler
var unfollowHandler *handler.UnfollowHandler

func init() {
	followHandler = &handler.FollowHandler{}
	unfollowHandler = &handler.UnfollowHandler{}
}

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
}

func setUpFollowHandler(m *tb.Message, bot *tb.Bot) int64 {
	userID := getUserID(m.Sender)
	steamID, err := getSteamID(m.Payload)

	log.Printf("M=setUpFollowHandler userID=%v steamID=%v\n", userID, steamID)

	if err != nil || len(steamID) != steamIDLength {
		bot.Send(m.Sender, "Invalid Param!")
		return -1
	}

	currNickname := GetPlayerCurrentNickname(steamID)

	return followHandler.FollowHandler(m, bot, steamID, currNickname, userID)
}

func setUpShowHandler(m *tb.Message, bot *tb.Bot) {
	userID := getUserID(m.Sender)

	handler.ShowHandler(m, bot, userID)
}

func setUpStopHandler(m *tb.Message, bot *tb.Bot) {
	userID := getUserID(m.Sender)

	handler.StopHandler(m, bot, userID)
}

func setUpStartHandler(m *tb.Message, bot *tb.Bot) {
	handler.StartHandler(m, bot)
}

func setUpUnfollowHandler(m *tb.Message, bot *tb.Bot) {
	steamID, err := getSteamID(m.Payload)

	log.Printf("M=setUpUnfollowHandler steamID=%v\n", steamID)

	if err != nil || len(steamID) != steamIDLength {
		bot.Send(m.Sender, "Invalid Param!")
		return
	}

	unfollowHandler.UnfollowHandler(m, bot, steamID)
}