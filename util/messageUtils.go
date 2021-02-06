package util

import (
	"fmt"
	"math/rand"
	"strings"
)

//GetNicknameChangedMessage returns the message when players change their nicknames
func GetNicknameChangedMessage(oldNickname, recentNickname, currNickname, steamID string) string {
	diffRecentNickname := ""
	if oldNickname != recentNickname {
		diffRecentNickname = ", recently playing as \"" + recentNickname + "\""
	}

	return "NICKNAME CHANGED: The user you followed as \"" + oldNickname + "\"" + diffRecentNickname + ", Steam Profile: " + SteamProfileURL + steamID + ", is now under the nickname \"" + currNickname + "\""
}

//GetBanMessage returns the message when a player get banned
func GetBanMessage(oldNickname, currNickname, steamID string, daysSinceLastBan int) string {
	// if the player hasn't changed nickname no reason to return redundant message
	changedNickPhrase := ""
	if oldNickname != currNickname {
		changedNickPhrase = ", now under the nickname " + strings.Replace(currNickname, "%", "", -1)
	}
	return "❌❌❌ BAN NEWS: The user you followed as " + strings.Replace(oldNickname, "%", "", -1) + changedNickPhrase + ", Steam Profile: " + SteamProfileURL + steamID + ", has just been BANNED! You won't be notified about this player anymore."
}

//GetFollowResponseMessage returns the message when a user follow a player
func GetFollowResponseMessage(oldNickname, currNickname string, followersCount int64, isVACBanned bool) string {

	var status string = "NOT banned (yet)."
	if isVACBanned {
		status = "BANNED (yay)."
	}

	if oldNickname != "" {

		messagePartOne := fmt.Sprintf("You already follow %v. ", currNickname)
		messagePartTwo := fmt.Sprintf("")
		messagePartThree := fmt.Sprintf("His current status is: %v.", status)

		if currNickname != oldNickname {
			messagePartTwo = fmt.Sprintf("He used to go by %v. ", oldNickname)
		}

		return messagePartOne + messagePartTwo + messagePartThree
	}

	message := fmt.Sprintf("Following player %v, status=%v", currNickname, status)

	if !isVACBanned {
		if followersCount > 0 {
			message += fmt.Sprintf(" Which is being followed by %v other users.", followersCount)
		} else {
			message += " You're the first to follow this player!"
		}
	} else {
		message += " You will NOT receive updates about this player since it's banned!"
	}
	return message
}

//GetStartResponse returns the message for the /start handler
func GetStartResponse(username string) string {
	return fmt.Sprintf("Hello %v, welcome to the Is VAC Banned Bot! EN/PT embaixo\n"+
		"These are the commands accepted by the bot: \n"+
		"/start : show you this very same message;\n\n"+
		"/stop : to disable any notifications you signed up for before;\n\n"+
		"/follow <argument>: follow a player, you'll be notified about nickname changes and VAC bans;\n"+
		"if you have previously disabled notifications, they will be enabled again.\n"+
		"<argument> can be either a Steam ID, Custom ID or the player's profile URL.\n\n"+
		"/unfollow <argument>: disables notifications about a specific player;\n"+
		"<argument> MUST be either a Steam ID or Custom ID!\n\n"+
		"/show : displays a list of every player you're following with its BAN status;\n\n"+
		"/summary : displays the summary of the players you follow with performance percentage;\n\n"+
		"Consider donating to help keep this service up and running: https://picpay.me/felipedeckert \n\n"+
		"________________________________________________\n"+
		"Esses são os comandos aceitos pelo bot: \n"+
		"/start : te mostra essa mesma mensagem;\n\n"+
		"/stop : desativa todas as notificações que você se inscreveu anteriormente;\n\n"+
		"/follow <argumento>: segue um jogador, você será notificado sobre todas trocas de nome e banimentos VAC;\n"+
		"se você anteriormente desabilitou notificações, elas serão ativadas novamente.\n"+
		"<argumento> pode ser Steam ID, Custom ID ou a URL do perfil.\n\n"+
		"/unfollow <argumento>: desabilita notificações de um jogador específico;\n"+
		"<argumento> tem que ser Steam ID ou Custom ID!\n\n"+
		"/show : mostra uma lista de todos os jkogadores seguidos e seus status de banimento;\n\n"+
		"/summary : mostra um resumo dos jogadores seguidos e porcentagem de banimento;\n\n"+
		"Considere doar para manter o serviço funcionando e com atualizações: https://picpay.me/felipedeckert", username)
}

//GetStopResponse returns the message for the /stop handler
func GetStopResponse() string {
	return "You will not be notified about any player anymore! Follow another player to start receiving news about all the players you followed."
}

//GetSummaryResponse returns the message for the /summary handler
func GetSummaryResponse(summary map[bool]int) string {
	messageEnd := ""

	if summary[false]+summary[true] == 0 {
		messageEnd = ", let's start tracking some suspects!"
	} else if summary[true] == 0 {
		messageEnd = fmt.Sprintf(", no one have been banned yet. Don't worry VAC will get them!")
	} else {
		messageEnd = fmt.Sprintf(", of which %v have been banned, a performance of %.2f%%. Keep up the good work!", summary[true], (float64(summary[true])/float64(summary[false]+summary[true]))*100)
	}

	return fmt.Sprintf("You follow %v players%v", summary[true]+summary[false], messageEnd)
}

//GetChooserResponse returns the message for the /chooser handler
func GetChooserResponse(i int) string {
	CS := 0

	if i == CS {
		return getCSPhrase()
	}

	return getValorantPhrase()
}

func getCSPhrase() string {
	res := rand.Intn(2)
	switch res {
	case 0:
		return "CS né, bicho?! Cês já viram a faquinha do Deck?!?"
	default:
		return "CS, mas o nosso Supremo tem que garantir 30+ kills!"
	}
}

func getValorantPhrase() string {
	res := rand.Intn(3)
	switch res {
	case 0:
		return "MAGISTRAL: bora ter aula com esse nosso OMEN, XNDão!"
	case 1:
		return "Valozinho né?! Sempre tem AQUELE ace do Vilão!"
	default:
		return "Bora Valo que hoje tem: clutch do BOSS!"
	}
}
