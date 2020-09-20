package main

import (
	"isvacbanned/job"
	"isvacbanned/service"
	"log"
	"os"

	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {

	//Receives a google drive spreadsheet ID as argument and updates vac ban status
	//sheetID := os.Args[1]
	//fmt.Printf("M=main spreadsheetID=%v step=2 \n", sheetID)
	//service.UpdatePlayersStatus(sheetID)

	startTelegramBot()
	go job.RunScheduler()

}

func startTelegramBot() {
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

	service.SetUpBot(webhook, token)

	log.Println("END")
}
