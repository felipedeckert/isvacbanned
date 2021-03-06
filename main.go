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
	var (
		port      = os.Getenv("PORT")
		publicURL = os.Getenv("PUBLIC_URL")
		token     = os.Getenv("TOKEN")
	)

	util.StartDatabase()

	webhook := &tb.Webhook{
		Listen:   ":" + port,
		Endpoint: &tb.WebhookEndpoint{PublicURL: publicURL},
	}

	go job.RunScheduler()

	service.SetUpBot(webhook, token)
}
