package main

import (
	"isvacbanned/job"
	"isvacbanned/service"
	"isvacbanned/util"
	"log"
	"os"

	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {

	//Receives a google drive spreadsheet ID as argument and updates vac ban status
	/*
		sheetID := os.Args[1]
		fmt.Printf("M=main spreadsheetID=%v step=2 \n", sheetID)
		service.UpdatePlayersStatus(sheetID)
	*/

	startTelegramBot()
	job.RunScheduler()

}

func startTelegramBot() {
	var (
		port      = os.Getenv("PORT")
		publicURL = os.Getenv("PUBLIC_URL")
		token     = os.Getenv("TOKEN")
	)

	if util.LOCAL {
		port = "3000"                                            //os.Getenv("PORT")
		publicURL = "https://is-vac-banned.herokuapp.com/"       //os.Getenv("PUBLIC_URL")
		token = "1262870496:AAG_XdC_OYONPVWeGAxInBmAGr2JfT8uOl0" //os.Getenv("TOKEN")
	}

	util.StartDatabase()

	webhook := &tb.Webhook{
		Listen:   ":" + port,
		Endpoint: &tb.WebhookEndpoint{PublicURL: publicURL},
	}

	go job.RunScheduler()

	service.SetUpBot(webhook, token)

	log.Println("END")
}
