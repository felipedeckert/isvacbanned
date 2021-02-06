package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"isvacbanned/util"
	"log"
	"net/http"
)

const noMatch int = 42

var Client HTTPClient

// PlayerSteamID is the Player Steam ID Struct returned by the steam API (userURL)
type PlayerSteamID struct {
	Response responseData `json:"response"`
}

type PlayerNickname struct {
	Response playerNicknameData `json:"response"`
}

type responseNicknameData struct {
	Personaname string `json:"personaname"`
}

type responseData struct {
	SteamId string `json:"steamId"`
	Success int    `json:"success"`
}

type playerData struct {
	SteamId          string `json:"SteamId"`
	CommunityBanned  bool   `json:"CommunityBanned"`
	VACBanned        bool   `json:"VACBanned"`
	NumberOfVACBans  int    `json:"NumberOfVACBans"`
	DaysSinceLastBan int    `json:"DaysSinceLastBan"`
	NumberOfGameBans int    `json:"NumberOfGameBans"`
	EconomyBan       string `json:"EconomyBan"`
}

type playerNicknameData struct {
	Players []responseNicknameData `json:"players"`
}

// Player represents a list of players BAN data
type Player struct {
	Players []playerData `json:"players"`
}

type playerService struct{}

type PlayerServiceInterface interface {
	GetPlayerStatus(steamID string) Player
	GetPlayerCurrentNickname(steamID string) string
}

var PlayerServiceClient PlayerServiceInterface = playerService{}

func init() {
	Client = &http.Client{}
}

func updatePlayersIfNeeded(players map[string]Player, spreadsheetID string) {
	log.Printf("M=updatePlayersIfNeeded")
	for idx, p := range players {
		data := p.Players[0]
		if data.VACBanned {
			log.Printf("M=updatePlayersIfNeeded VAC " + data.SteamId)
			UpdateVACBanStatus(idx, data.DaysSinceLastBan, spreadsheetID)
		}
	}
}

// UnmarshalPlayerByID returns a player and its data obtained from Steam API
func unmarshalPlayer(jsonInput []byte) Player {
	player := Player{}

	err := json.Unmarshal(jsonInput, &player)

	if err != nil {
		panic(err)
	}
	return player
}

func unmarshalSteamID(str []byte) (string, error) {
	playerID := PlayerSteamID{}

	err := json.Unmarshal(str, &playerID)

	if err != nil {
		panic(err)
	}

	if playerID.Response.Success == noMatch {
		return "", errors.New("User not found")
	}

	return playerID.Response.SteamId, nil
}

/*
func getAllPlayersStatus(userSteamID map[string]string) map[string]Player {
	log.Printf("M=getAllPlayersStatus")
	players := make(map[string]Player)
	for idx, value := range userSteamID {
		players[idx] = GetPlayerStatus(value)
	}
	return players
}
*/

// GetPlayerStatus receives a steamID and returns its player ban status
func (p playerService) GetPlayerStatus(steamID string) Player {
	url := util.GetVACBanURL(steamID)
	//todo change to debug level
	//log.Printf("M=getPlayerStatus url=%v\n", url)
	resp, err := Client.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	result, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	return unmarshalPlayer(result)
}

func getPlayerSteamID(playerName string) (string, error) {
	url := util.GetNicknameURL(playerName)
	log.Printf("M=getPlayerSteamID playerName=%v\n", playerName)
	resp, err := Client.Get(url)

	if err != nil {
		log.Printf("M=getPlayerSteamID step=get err=%s\n", err)
		return "", errors.New("unable to get player")
	}

	result, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Printf("M=getPlayerSteamID step=parse err=%s\n", err)
		return "", errors.New("unable to parse player")
	}

	res, err := unmarshalSteamID(result)

	if err != nil {
		log.Printf("M=getPlayerSteamID step=unmarshal err=%s\n", err)
		return "", err
	}

	return res, nil
}

// GetPlayerCurrentNickname gets the player identified by steamID current nickname
func (p playerService) GetPlayerCurrentNickname(steamID string) string {
	url := util.GetPlayerSummaryURL(steamID)
	resp, err := Client.Get(url)

	if err != nil {
		log.Printf("M=GetPlayerCurrentNickname err=%s\n", err)
		log.Fatal(err)
	}

	result, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	playerNickname := PlayerNickname{}

	err = json.Unmarshal(result, &playerNickname)

	if err != nil {
		panic(err)
	}

	if len(playerNickname.Response.Players) == 0 {
		log.Printf("M=GetPlayerCurrentNickname status=404 steamID=%v\n", steamID)
		return ""
	} /*else {
		//todo change to debug level
		log.Printf("M=GetPlayerCurrentNickname status=200 steamID=%v\n", steamID)
	}*/

	return playerNickname.Response.Players[0].Personaname

}

/*
//UpdatePlayersStatus updates players status on given spreadsheet
func UpdatePlayersStatus(spreadsheetID string) {
	log.Printf("M=UpdatePlayersStatus spreadsheetID=%v\n", spreadsheetID)

	userSteamID := GetSteamIDs(spreadsheetID)

	players := getAllPlayersStatus(userSteamID)

	updatePlayersIfNeeded(players, spreadsheetID)
}
*/
