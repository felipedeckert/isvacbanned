package service

import (
	"fmt"
	"isvacbanned/model"
	"log"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

const steamIDLength = 17
const playerBanned = "This player is VAC banned!"
const playerNotBanned = "This player is NOT VAC banned!"

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

	bot.Handle("/show", func(m *tb.Message) { showHandler(m, bot) })
}

func verifyHandler(m *tb.Message, bot *tb.Bot) {
	steamID, err := getSteamID(m.Payload)

	if err != nil || len(steamID) != steamIDLength {
		bot.Send(m.Sender, "Invalid Param!")
		return
	}

	player := GetPlayerStatus(steamID)

	if len(player.Players) == 0 {
		bot.Send(m.Sender, "Player not found!")
		return
	}

	isVACBanned := player.Players[0].VACBanned

	fmt.Printf("M=verifyHandler player=%v isVACBanned=%v\n", steamID, isVACBanned)

	message := getBanResponse(isVACBanned)

	bot.Send(m.Sender, message)
}

func followHandler(m *tb.Message, bot *tb.Bot) {
	userID := getUserID(m.Sender)
	steamID, err := getSteamID(m.Payload)

	fmt.Printf("M=followHandler payload=%v chatID=%v\n", m.Payload, m.Chat.ID)
	if err != nil || len(steamID) != steamIDLength {
		bot.Send(m.Sender, "Invalid Param!")
		return
	}

	currNickname := GetPlayerCurrentNickname(steamID)

	followersCount, err := model.GetFollowerCountBySteamID(steamID)

	if err != nil {
		panic(err)
	}

	model.FollowSteamUser(m.Chat.ID, steamID, currNickname, userID)

	message := getFollowResponse(currNickname, followersCount)

	bot.Send(m.Sender, message)
}

func showHandler(m *tb.Message, bot *tb.Bot) {
	userID := getUserID(m.Sender)
	followedUsers := model.GetUsersFollowed(userID)

	message := getShowResponse(followedUsers)

	bot.Send(m.Sender, message)
}

func getShowResponse(followedUsers []model.UsersFollowed) string {
	if len(followedUsers) == 0 {
		return "You're not following any player yet!"
	}

	return "You're following these users: \n" + getPlayersAndStatusAsShoppingList(followedUsers)
}

func getPlayersAndStatusAsShoppingList(followedUsers []model.UsersFollowed) string {
	var str strings.Builder

	for _, user := range followedUsers {
		status := "NOT BANNED"
		if user.IsCompleted {
			status = "BANNED"
		}
		str.WriteString(user.OldNickname + " : " + status + ",\n")
	}

	return str.String()
}

func getFollowResponse(currNickname string, followersCount int64) string {

	message := fmt.Sprintf("Following player %v, ", currNickname)

	if followersCount > 0 {
		message += fmt.Sprintf("which is being followes by %v other users.", followersCount)
	} else {
		message += "you're the first to follow this player!"
	}

	return message
}

func getBanResponse(isVACBanned bool) string {
	if isVACBanned {
		return playerBanned
	}
	return playerNotBanned
}
