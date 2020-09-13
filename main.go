package main

import (
	"isvacbanned/service"
	"log"
	"os"

	tb "gopkg.in/tucnak/telebot.v2"
)

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

func setUpBotHandlers(bot *tb.Bot) {
	bot.Handle("/hello", func(m *tb.Message) {
		bot.Send(m.Sender, "Hi!")
	})

	bot.Handle("/verify", func(m *tb.Message) {
		service.UnmarshalPlayer(m.Payload)

		bot.Send(m.Sender, m.Payload)
	})

}

func updatePlayersStatus() {
	userSteamID := service.GetSteamIDs()

	players := service.GetAllPlayersStatuses(userSteamID)

	service.UpdatePlayersIfNeeded(players)
}
