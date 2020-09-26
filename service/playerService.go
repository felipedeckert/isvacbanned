package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const valveKey = "DD5F4C5D083B1C9F7AB2CCAC76124DEC"
const vacBanURL = "http://api.steampowered.com/ISteamUser/GetPlayerBans/v1/?key="
const userURL = "http://api.steampowered.com/ISteamUser/ResolveVanityURL/v0001/?key="
const playerSummaryURL = "http://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key="
const userParamKey = "&vanityurl="
const steamIDParamKey = "&steamids="

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

func init() {
	Client = &http.Client{}
}

func updatePlayersIfNeeded(players map[string]Player, spreadsheetID string) {
	fmt.Println("M=updatePlayersIfNeeded")
	for idx, p := range players {
		data := p.Players[0]
		if data.VACBanned {
			fmt.Println("M=updatePlayersIfNeeded VAC " + data.SteamId)
			UpdateVACBanStatus(idx, data.DaysSinceLastBan, spreadsheetID)
		}
	}
}

// UnmarshalPlayerByID returns a player and its data obtained from Steam API
func unmarshalPlayerByID(jsonInput []byte) Player {
	player := Player{}

	err := json.Unmarshal(jsonInput, &player)

	if err != nil {
		panic(err)
	}
	return player
}

func unmarshalPlayerByName(customID string) string {
	log.Printf("M=UnmarshalPlayerByName customID%v=", customID)
	playerID := PlayerSteamID{}
	str, err := getPlayerSteamID(customID)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(str, &playerID)

	if err != nil {
		panic(err)
	}

	return playerID.Response.SteamId
}

func getAllPlayersStatus(userSteamID map[string]string) map[string]Player {
	fmt.Println("M=getAllPlayersStatus")
	players := make(map[string]Player)
	for idx, value := range userSteamID {

		players[idx] = GetPlayerStatus(value)
		fmt.Println(value)
	}
	return players
}

func GetPlayerStatus(steamID string) Player {
	url := buildGetURL(steamID)
	log.Printf("M=getPlayerStatus url=%v\n", url)
	resp, err := Client.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	result, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	return unmarshalPlayerByID(result)
}

func getPlayerSteamID(playerName string) ([]byte, error) {
	url := buildGetUserURL(playerName)
	log.Printf("M=getPlayerSteamID url=%v\n" + url)
	resp, err := Client.Get(url)

	if err != nil {
		log.Printf("M=getPlayerSteamID err=%s\n", err)
		return nil, err
	}

	result, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	return result, nil
}

// GetPlayerCurrentNickname gets the player identified by steamID current nickname
func GetPlayerCurrentNickname(steamID string) string {
	url := buildGetPlayerSummaryURL(steamID)
	log.Printf("M=getPLayerCurrentNickname url=%v\n", url)
	resp, err := Client.Get(url)

	if err != nil {
		log.Printf("M=getPlayerSteamID err=%s\n", err)
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

	return playerNickname.Response.Players[0].Personaname

}

func buildGetPlayerSummaryURL(steamID string) string {
	fmt.Printf("M=buildGetPlayerSummaryURL SteamID=%s\n", steamID)
	return playerSummaryURL + valveKey + steamIDParamKey + steamID
}

func buildGetUserURL(userName string) string {
	fmt.Printf("M=buildGetUserURL userName=%s\n", userName)
	return userURL + valveKey + userParamKey + userName
}

func buildGetURL(steamID string) string {
	fmt.Printf("M=buildGetURL SteamID=%s\n", steamID)
	return vacBanURL + valveKey + steamIDParamKey + steamID
}

//UpdatePlayersStatus updates players status on given spreadsheet
func UpdatePlayersStatus(spreadsheetID string) {
	fmt.Printf("M=UpdatePlayersStatus spreadsheetID=%v\n", spreadsheetID)

	userSteamID := GetSteamIDs(spreadsheetID)

	players := getAllPlayersStatus(userSteamID)

	updatePlayersIfNeeded(players, spreadsheetID)
}
