package handler

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
)

//StopHandler handles show requests
func StartHandler(m *tb.Message, bot *tb.Bot) {
	message := getStartResponse(m.Sender.FirstName)

	bot.Send(m.Sender, message)
}

func getStartResponse(username string) string {
	return fmt.Sprintf("Hello %v, welcome to the Is VAC Banned Bot! \n"+
		"These are the commands accepted by the bot: \n"+
		"/start : show you this very same message;\n\n"+
		"/stop : to disable any notifications you signed up for before;\n\n"+
		"/follow <argument>: follow a player, you'll be notified about nickname changes and VAC bans;\n"+
		"if you have previously disabled notifications, they will be enabled again.\n"+
		"<argument> can be either a Steam ID or the player's profile URL.\n\n"+
		"/unfollow <argument>: disables notifications about a specific player;\n"+
		"<argument> MUST be a Steam ID!\n\n"+
		"/show : displays a list of every player you're following with its BAN status;\n\n"+
		"Consider donating to help keep this service up and running: linkToPaypal.com", username)
}
