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
const paramKey = "&steamids="

type playerData struct {
	SteamId          string `json:"SteamId"`
	CommunityBanned  bool   `json:"CommunityBanned"`
	VACBanned        bool   `json:"VACBanned"`
	NumberOfVACBans  int    `json:"NumberOfVACBans"`
	DaysSinceLastBan int    `json:"DaysSinceLastBan"`
	NumberOfGameBans int    `json:"NumberOfGameBans"`
	EconomyBan       string `json:"EconomyBan"`
}

type player struct {
	Players []playerData `json:"players"`
}

func UpdatePlayersIfNeeded(players map[string]player) {
	for idx, p := range players {
		data := p.Players[0]
		if data.VACBanned {
			UpdateVACBanStatus(idx, data.DaysSinceLastBan)
		}
	}
}

func GetAllPlayersStatuses(userSteamID map[string]string) map[string]player {
	players := make(map[string]player)
	for idx, value := range userSteamID {
		player := player{}
		str, err := getPlayerStatus(value)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", str)
		err = json.Unmarshal(str, &player)

		if err != nil {
			panic(err)
		}
		players[idx] = player
		fmt.Println(player)
	}
	return players
}

func getPlayerStatus(steamID string) ([]byte, error) {
	url := buildGetURL(steamID)
	fmt.Println(url)
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

func buildGetURL(steamID string) string {
	fmt.Printf("SteamID=%s\n", steamID)
	return vacBanURL + valveKey + paramKey + steamID
}
