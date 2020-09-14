package main

import (
	"fmt"
	"isvacbanned/service"
	"log"
	"os"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

const steamIDLength = 17

func main() {

	var (
		port      = os.Getenv("PORT")
		publicURL = os.Getenv("PUBLIC_URL")
		token     = os.Getenv("TOKEN")
	)

	local := false
	if local {
		port = "3030"                                            //os.Getenv("PORT")
		publicURL = "https://is-vac-banned.herokuapp.com/"       //os.Getenv("PUBLIC_URL")
		token = "1324910657:AAFSlJn6TD9EeYNn35MEo-YphYlhYhqc_do" //os.Getenv("TOKEN")
	}

	webhook := &tb.Webhook{
		Listen:   ":" + port,
		Endpoint: &tb.WebhookEndpoint{PublicURL: publicURL},
	}

	setUpBot(webhook, token)

	log.Println("END")
}

func setUpBot(webhook *tb.Webhook, token string) {
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

func getSteamID(input string) string {
	argument := input
	if strings.Contains(input, "id") {

		argument = strings.Split(input, "/")[2]

	} else if strings.Contains(input, "profile") {

		argument = strings.Split(input, "/")[2]

	}
	log.Printf("M=getSteamID input=%v argument=%v\n", input, argument)
	return argument
}

func verifyHandler(m *tb.Message, bot *tb.Bot) {
	argument := getSteamID(m.Payload)

	if len(argument) != steamIDLength {
		bot.Send(m.Sender, "Invalid Steam ID!")
		return
	}

	player := service.UnmarshalPlayerByID(argument)

	if len(player.Players) == 0 {
		bot.Send(m.Sender, "Player not found!")
		return
	}

	isVACBanned := player.Players[0].VACBanned

	fmt.Printf("M=verifyHandler player=%v isVACBanned=%v\n", argument, isVACBanned)

	result := getResponse(isVACBanned)

	bot.Send(m.Sender, result)
}

func setUpBotHandlers(bot *tb.Bot) {
	bot.Handle("/hello", func(m *tb.Message) {
		bot.Send(m.Sender, "Hi!")
	})

	bot.Handle("/verify", func(m *tb.Message) { verifyHandler(m, bot) })

}

func getResponse(isVACBanned bool) string {
	if isVACBanned {
		return "This player is VAC banned!"
	}

	return "This player is NOT VAC banned!"
}

func updatePlayersStatus() {
	userSteamID := service.GetSteamIDs()

	players := service.GetAllPlayersStatuses(userSteamID)

	service.UpdatePlayersIfNeeded(players)
}
