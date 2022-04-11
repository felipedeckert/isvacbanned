package job

import (
	"isvacbanned/messenger"
	"isvacbanned/model"
	"isvacbanned/service"
	"isvacbanned/util"
	"log"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

// RunScheduler sets up the scheduler to check users status
func RunScheduler() {
	log.Printf("M=RunScheduler\n")
	scheduler := gocron.NewScheduler(time.UTC)

	scheduler.Every(60).Minutes().Do(checkFollowedUsersBan)
	scheduler.Every(10).Minutes().Do(checkFollowedUsersNickname)

	scheduler.StartAsync()
}

func checkFollowedUsersBan() {
	log.Print("M=checkFollowedUsersBan step=start")
	currTime := time.Now()
	usersIncomplete := model.FollowModelClient.GetAllIncompleteFollowedUsers()
	var wg sync.WaitGroup

	for chatID, steamIDList := range usersIncomplete {
		for _, users := range steamIDList {
			wg.Add(1)
			go validateBanStatusAndSendMessage(users, chatID, &wg)
		}
	}
	wg.Wait()
	elapsedTime := time.Since(currTime)
	log.Printf("M=checkFollownUsersBan step=end et=%dms", int64(elapsedTime/time.Millisecond))
}

func checkFollowedUsersNickname() {
	log.Print("M=checkFollowedUsersNickname step=start")
	currTime := time.Now()
	usersIncomplete := model.FollowModelClient.GetAllIncompleteFollowedUsers()

	var wg sync.WaitGroup

	for chatID, steamIDList := range usersIncomplete {
		for _, users := range steamIDList {
			wg.Add(1)
			go validateNicknameAndSendMessage(users, chatID, &wg)
		}
	}
	wg.Wait()
	elapsedTime := time.Since(currTime)
	log.Printf("M=checkFollowedUsersNickname step=end et=%dms", int64(elapsedTime/time.Millisecond))
}

func validateBanStatusAndSendMessage(user model.UsersFollowed, chatID int64, wg *sync.WaitGroup) {
	defer wg.Done()
	idsToUpdate := make([]int64, 0)
	player := service.PlayerServiceClient.GetPlayerStatus(user.SteamID)

	if len(player.Players) > 0 {
		playerData := player.Players[0]
		if playerData.VACBanned {
			actualNickname := service.PlayerServiceClient.GetPlayerCurrentNickname(user.SteamID)
			if actualNickname != "" {
				log.Printf("M=validateBanStatusAndSendMessage steamID=%v status=banned\n", user.SteamID)
				messenger.MessengerClient.SendMessageToUser(util.GetBanMessage(user.OldNickname, actualNickname, user.SteamID, playerData.DaysSinceLastBan), chatID)
				idsToUpdate = append(idsToUpdate, user.ID)
			}
		}
	}

	if len(idsToUpdate) > 0 {
		//Once a player status is set to completed, this player will not be returned in the GetAllIncompleteFollowedUsers query
		model.FollowModelClient.SetFollowedUserToCompleted(idsToUpdate)
	}
}

func validateNicknameAndSendMessage(user model.UsersFollowed, chatID int64, wg *sync.WaitGroup) {
	defer wg.Done()
	actualNickname := service.PlayerServiceClient.GetPlayerCurrentNickname(user.SteamID)

	//sometimes the request fails, and it comes empty, hence this validation
	if actualNickname == "" {
		log.Printf("M=validateNicknameAndSendMessage L=E steamID=%v status=emptyCurrentNickname\n", user.SteamID)
	} else if user.CurrNickname != actualNickname {
		log.Printf("M=validateNicknameAndSendMessage L=I steamID=%v status=changedNickname\n", user.SteamID)
		model.FollowModelClient.SetCurrNickname(user.ID, actualNickname)
		messenger.MessengerClient.SendMessageToUser(util.GetNicknameChangedMessage(user.OldNickname, user.CurrNickname, actualNickname, user.SteamID), chatID)
	}
}
