package main

import (
	"isvacbanned/job"
	"isvacbanned/service"
	"isvacbanned/util"
	"os"

	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	startTelegramBot()
	job.RunScheduler()
}

func startTelegramBot() {
	port      	  := os.Getenv("PORT")
	publicURL 	  := os.Getenv("PUBLIC_URL")
	telegramToken := os.Getenv("TOKEN")

	util.StartDatabase()

	webhook := &tb.Webhook{
		Listen:   ":" + port,
		Endpoint: &tb.WebhookEndpoint{PublicURL: publicURL},
	}

	go job.RunScheduler()

	service.SetUpBot(webhook, os.Getenv(telegramToken))
}
