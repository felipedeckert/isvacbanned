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
const userParamKey = "&vanityurl="
const steamIDParamKey = "&steamids="

type PlayerSteamID struct {
	Response responseData `json:"response"`
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

type Player struct {
	Players []playerData `json:"players"`
}

func UpdatePlayersIfNeeded(players map[string]Player) {
	for idx, p := range players {
		data := p.Players[0]
		if data.VACBanned {
			UpdateVACBanStatus(idx, data.DaysSinceLastBan)
		}
	}
}

func UnmarshalPlayerByID(value string) Player {
	player := Player{}
	str, err := getPlayerStatus(value)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", str)
	err = json.Unmarshal(str, &player)

	if err != nil {
		panic(err)
	}
	return player
}

func UnmarshalPlayerByName(value string) string {
	playerID := PlayerSteamID{}
	str, err := getPlayerSteamID(value)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(str, &playerID)

	if err != nil {
		panic(err)
	}

	return playerID.Response.SteamId
}

func GetAllPlayersStatuses(userSteamID map[string]string) map[string]Player {
	players := make(map[string]Player)
	for idx, value := range userSteamID {

		players[idx] = UnmarshalPlayerByID(value)
		fmt.Println(value)
	}
	return players
}

func getPlayerStatus(steamID string) ([]byte, error) {
	url := buildGetURL(steamID)
	fmt.Println("M=getPlayerStatus url=" + url)
	resp, err := http.Get(url)

	if err != nil {
		fmt.Printf("M=getPlayerStatus err=%s\n", err)
		return nil, err
	}

	result, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	return result, nil
}

func getPlayerSteamID(playerName string) ([]byte, error) {
	url := buildGetUserURL(playerName)
	fmt.Println("M=getPlayerSteamID url=" + url)
	resp, err := http.Get(url)

	if err != nil {
		fmt.Printf("M=getPlayerSteamID err=%s\n", err)
		return nil, err
	}

	result, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	return result, nil
}

func buildGetUserURL(userName string) string {
	fmt.Printf("M=buildGetUserURL userName=%s\n", userName)
	return userURL + valveKey + userParamKey + userName
}

func buildGetURL(steamID string) string {
	fmt.Printf("M=buildGetURL SteamID=%s\n", steamID)
	return vacBanURL + valveKey + userParamKey + steamID
}
