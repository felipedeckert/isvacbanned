package service

import (
	"isvacbanned/handler"
	"isvacbanned/util"
	"log"
	"strconv"

	tb "gopkg.in/tucnak/telebot.v2"
)

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
	userID := getUserID(m.Chat)
	steamID, err := getSteamID(m.Payload)

	log.Printf("M=setUpFollowHandler payload=%v userID=%v steamID=%v\n", m.Payload, userID, steamID)

	if err != nil {
		bot.Send(m.Chat, err.Error())
		return -1
	} else if len(steamID) != util.SteamIDLength || !isNumeric(steamID) {
		bot.Send(m.Chat, "Invalid Param!")
		return -1
	}

	currNickname := GetPlayerCurrentNickname(steamID)

	player := GetPlayerStatus(steamID)
	playerData := player.Players[0]

	return followHandler.FollowHandler(m, bot, steamID, currNickname, userID, playerData.VACBanned)
}

func setUpShowHandler(m *tb.Message, bot *tb.Bot) {
	userID := getUserID(m.Chat)

	handler.ShowHandler(m, bot, userID)
}

func setUpStopHandler(m *tb.Message, bot *tb.Bot) {
	userID := getUserID(m.Chat)

	handler.StopHandler(m, bot, userID)
}

func setUpStartHandler(m *tb.Message, bot *tb.Bot) {
	handler.StartHandler(m, bot)
}

func setUpUnfollowHandler(m *tb.Message, bot *tb.Bot) {
	userID := getUserID(m.Chat)
	steamID, err := getSteamID(m.Payload)

	log.Printf("M=setUpUnfollowHandler steamID=%v\n", steamID)

	if err != nil {
		bot.Send(m.Chat, err.Error())
		return
	} else if len(steamID) != util.SteamIDLength || !isNumeric(steamID) {
		bot.Send(m.Chat, "Invalid Param!")
		return
	}

	unfollowHandler.UnfollowHandler(m, bot, steamID, userID)
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
