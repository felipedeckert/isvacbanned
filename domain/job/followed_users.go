package job

import (
	"context"
	"isvacbanned/domain/entities"
	"isvacbanned/util"
	"log"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

//go:generate moq -stub -pkg mocks -out mocks/followed_users_job_steam_gateway.go . FollowedUsersJobSteamGateway
//go:generate moq -stub -pkg mocks -out mocks/followed_users_job_telegram_gateway.go . FollowedUsersJobTelegramGateway
//go:generate moq -stub -pkg mocks -out mocks/followed_users_job_follow_persistence_gateway.go . FollowedUsersJobFollowPersistenceGateway

const steamMaxInputLength = 100

type FollowedUsersJobFollowPersistenceGateway interface {
	GetAllIncompleteFollowedUsers(ctx context.Context) (map[int64][]entities.UsersFollowed, error)
	SetFollowedUserToCompleted(ctx context.Context, id []int64)
	SetCurrNickname(ctx context.Context, userID int64, sanitizedActualNickname string) error
}

type FollowedUsersJobSteamGateway interface {
	GetPlayersStatus(steamIDs ...string) (entities.Player, error)
	GetPlayersCurrentNicknames(steamIDs ...string) (entities.PlayerInfo, error)
}

type FollowedUsersJobTelegramGateway interface {
	SendMessageToUser(message string, chatID int64)
}

type FollowedUsersJobUseCase struct {
	followPersistenceGateway FollowedUsersJobFollowPersistenceGateway
	steamGateway FollowedUsersJobSteamGateway
	telegramGateway FollowedUsersJobTelegramGateway
}

func NewFollowedUsersJobUseCase(
	followPersistenceGateway FollowedUsersJobFollowPersistenceGateway,
	steamGateway FollowedUsersJobSteamGateway,
	telegramGateway FollowedUsersJobTelegramGateway) *FollowedUsersJobUseCase {

	return &FollowedUsersJobUseCase{
		followPersistenceGateway: followPersistenceGateway,
		steamGateway:             steamGateway,
		telegramGateway:          telegramGateway,
	}
}

// RunScheduler sets up the scheduler to check users status
func (uc *FollowedUsersJobUseCase) RunScheduler(ctx context.Context, checkBanIntervalInMinutes, checkNicknameIntervalInMinutes uint64) {
	log.Printf("M=RunScheduler\n")
	scheduler := gocron.NewScheduler(time.UTC)

	scheduler.Every(checkBanIntervalInMinutes).Minutes().Do(uc.CheckFollowedUsersBan, ctx)
	scheduler.Every(checkNicknameIntervalInMinutes).Minutes().Do(uc.CheckFollowedUsersNickname, ctx)

	scheduler.StartAsync()
}

func (uc *FollowedUsersJobUseCase) CheckFollowedUsersBan(ctx context.Context) {
	log.Print("M=checkFollowedUsersBan step=start")
	currTime := time.Now()
	
	usersIncomplete, err := uc.followPersistenceGateway.GetAllIncompleteFollowedUsers(ctx)
	if err != nil {
		log.Printf("M=checkFollowedUsersBan error fetching users err:%s\n", err.Error())
	}
	
	var wg sync.WaitGroup
	var userGroup []entities.UsersFollowed
	for chatID, steamIDList := range usersIncomplete {
		for _, users := range steamIDList {
			userGroup = append(userGroup, users)
			if len(userGroup) % steamMaxInputLength == 0 {
				wg.Add(1)
				userGroupCopy := make([]entities.UsersFollowed, steamMaxInputLength)
				copy(userGroupCopy, userGroup)
				userGroup = make([]entities.UsersFollowed, 0)
				go func() {
					err = uc.validateBanStatusAndSendMessage(ctx, chatID, &wg, userGroupCopy...)
					if err != nil {
						log.Printf(`M=checkFollowedUsersBan step=1 L=E error validating status err:%s`, err.Error())
					}
				}()
			}
		}
		if len(userGroup) > 0 {
			wg.Add(1)
			userGroupCopy := make([]entities.UsersFollowed, len(userGroup))
			copy(userGroupCopy, userGroup)
			userGroup = make([]entities.UsersFollowed, 0)
			go func() {
				err = uc.validateBanStatusAndSendMessage(ctx, chatID, &wg, userGroupCopy...)
				if err != nil {
					log.Printf(`M=checkFollowedUsersBan step=2 L=E error validating status err:%s`, err.Error())
				}
			}()
		}
	}
	wg.Wait()
	elapsedTime := time.Since(currTime)
	log.Printf("M=checkFollownUsersBan step=end et=%dms", int64(elapsedTime/time.Millisecond))
}

func (uc *FollowedUsersJobUseCase) CheckFollowedUsersNickname(ctx context.Context) {
	log.Print("M=checkFollowedUsersNickname step=start")
	currTime := time.Now()
	usersIncomplete, err := uc.followPersistenceGateway.GetAllIncompleteFollowedUsers(ctx)
	if err != nil {
		log.Printf("M=checkFollowedUsersNickname error fetching users err:%s\n", err.Error())
	}
	
	var wg sync.WaitGroup
	var userGroup []entities.UsersFollowed
	for chatID, steamIDList := range usersIncomplete {
		for _, users := range steamIDList {
			userGroup = append(userGroup, users)
			if len(userGroup) % 100 == 0 {
				userGroupCopy := make([]entities.UsersFollowed, 100)
				copy(userGroupCopy, userGroup)
				userGroup = make([]entities.UsersFollowed, 0)
				wg.Add(1)
				go func() {
					err = uc.validateNicknameAndSendMessage(ctx, chatID, &wg, userGroupCopy...)
					if err != nil {
						log.Printf(`M=checkFollowedUsersNickname step=1 L=E error validating nickname err:%s`, err.Error())
					}
				}()
			}
		}
		if len(userGroup) > 0 {
			wg.Add(1)
			userGroupCopy := make([]entities.UsersFollowed, len(userGroup))
			copy(userGroupCopy, userGroup)
			userGroup = make([]entities.UsersFollowed, 0)
			go func() {
				err = uc.validateNicknameAndSendMessage(ctx, chatID, &wg, userGroupCopy...)
				if err != nil {
					log.Printf(`M=checkFollowedUsersNickname step=2 L=E error validating nickname err:%s`, err.Error())
				}
			}()
		}
	}

	wg.Wait()
	elapsedTime := time.Since(currTime)
	log.Printf("M=checkFollowedUsersNickname step=end et=%dms", int64(elapsedTime/time.Millisecond))
}

func (uc *FollowedUsersJobUseCase) validateBanStatusAndSendMessage(ctx context.Context, chatID int64, wg *sync.WaitGroup, userGroup ...entities.UsersFollowed) error {
	defer wg.Done()
	idsToUpdate := make([]int64, 0)
	steamIDs := make([]string, 0)
	users := make(map[string]entities.UsersFollowed)

	for _, user := range userGroup {
		steamIDs = append(steamIDs, user.SteamID)
		users[user.SteamID] = user
	}
	playersStatusInfo, err := uc.steamGateway.GetPlayersStatus(steamIDs...)
	if err != nil {
		return err
	}

	if len(playersStatusInfo.Players) > 0 {
		playersInfo, err := uc.steamGateway.GetPlayersCurrentNicknames(steamIDs...)
		if err != nil{
			return err
		}
		playersInfoMap := make(map[string]entities.ResponseNicknameData)
		for _, player := range playersInfo.Response.Players {
			playersInfoMap[player.SteamID] = player
		}

		for _, player := range playersStatusInfo.Players {
			if player.VACBanned {
				if playersInfoMap[player.SteamId].PersonaName != "" {
					log.Printf("M=validateBanStatusAndSendMessage steamID=%v status=banned\n", player.SteamId)
					message := util.GetBanMessage(users[player.SteamId].OldNickname, playersInfoMap[player.SteamId].PersonaName, player.SteamId, player.DaysSinceLastBan)
					uc.telegramGateway.SendMessageToUser(message, chatID)
					idsToUpdate = append(idsToUpdate, users[player.SteamId].ID)
				}
			}
		}
	}

	if len(idsToUpdate) > 0 {
		//Once a player status is set to completed, this player will not be returned in the GetAllIncompleteFollowedUsers query
		uc.followPersistenceGateway.SetFollowedUserToCompleted(ctx, idsToUpdate)
		idsToUpdate = make([]int64, 0)
	}

	return nil
}

func (uc *FollowedUsersJobUseCase) validateNicknameAndSendMessage(ctx context.Context, chatID int64, wg *sync.WaitGroup, userGroup ...entities.UsersFollowed) error {
	defer wg.Done()
	steamIDs := make([]string, 0)
	users := make(map[string]entities.UsersFollowed)

	for _, user := range userGroup {
		steamIDs = append(steamIDs, user.SteamID)
		users[user.SteamID] = user
	}

	playersInfo, err := uc.steamGateway.GetPlayersCurrentNicknames(steamIDs...)
	if err != nil{
		return err
	}

	for _, player := range playersInfo.Response.Players {
		//sometimes the request fails, and it comes empty, hence this validation
		if player.PersonaName == "" {
			log.Printf("M=validateNicknameAndSendMessage L=E steamID=%v status=emptyCurrentNickname\n", player.SteamID)
		} else if users[player.SteamID].CurrNickname != player.PersonaName {
			log.Printf("M=validateNicknameAndSendMessage L=I steamID=%v status=changedNickname\n", player.SteamID)
			err = uc.followPersistenceGateway.SetCurrNickname(ctx, users[player.SteamID].ID, player.PersonaName)
			if err != nil {
				log.Printf("M=validateNicknameAndSendMessage L=E err:%s\n", err.Error())
			}
			message := util.GetNicknameChangedMessage(users[player.SteamID].OldNickname, users[player.SteamID].CurrNickname, player.PersonaName, player.SteamID)
			uc.telegramGateway.SendMessageToUser(message, chatID)
		}
	}

	return nil
}
