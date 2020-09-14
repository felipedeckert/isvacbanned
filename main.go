package main

import (
	"errors"
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

func getArgumentFromURL(url string) (string, error) {
	splittedInput := strings.Split(url, "/")
	if len(splittedInput) == 0 {
		return "", errors.New("Invalid URL")
	}
	return splittedInput[len(splittedInput)-1], nil
}

func getSteamID(url string) (string, error) {
	steamID := url
	var err error
	var customID string
	if strings.Contains(url, "id") {
		customID, err = getArgumentFromURL(url)
		steamID = service.UnmarshalPlayerByName(customID)
	} else if strings.Contains(url, "profile") {
		steamID, err = getArgumentFromURL(url)
	}

	if err != nil {
		log.Printf("M=getSteamID input=%v\n", url)

		return "", err
	}

	log.Printf("M=getSteamID input=%v argument=%v\n", url, steamID)

	return steamID, nil
}

func verifyHandler(m *tb.Message, bot *tb.Bot) {
	argument, err := getSteamID(m.Payload)

	if err != nil || len(argument) != steamIDLength {
		bot.Send(m.Sender, "Invalid Param!")
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
