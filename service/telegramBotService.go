package service

import (
	"fmt"
	"isvacbanned/model"
	"log"

	tb "gopkg.in/tucnak/telebot.v2"
)

const steamIDLength = 17

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
	bot.Handle("/verify", func(m *tb.Message) { verifyHandler(m, bot) })

	bot.Handle("/follow", func(m *tb.Message) { followHandler(m, bot) })
}

func verifyHandler(m *tb.Message, bot *tb.Bot) {
	argument, err := getSteamID(m.Payload)

	if err != nil || len(argument) != steamIDLength {
		bot.Send(m.Sender, "Invalid Param!")
		return
	}

	player := UnmarshalPlayerByID(argument)

	if len(player.Players) == 0 {
		bot.Send(m.Sender, "Player not found!")
		return
	}

	isVACBanned := player.Players[0].VACBanned

	fmt.Printf("M=verifyHandler player=%v isVACBanned=%v\n", argument, isVACBanned)

	result := getResponse(isVACBanned)

	bot.Send(m.Sender, result)
}

func followHandler(m *tb.Message, bot *tb.Bot) {
	steamID, err := getSteamID(m.Payload)

	fmt.Printf("M=followHandler payload=%v chatID=%v\n", m.Payload, m.Chat.ID)
	if err != nil || len(steamID) != steamIDLength {
		bot.Send(m.Sender, "Invalid Param!")
		return
	}

	model.FollowSteamUser(m.Chat.ID, steamID)

	bot.Send(m.Sender, "Following user "+steamID)
}
