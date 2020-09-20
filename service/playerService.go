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

// PlayerSteamID is the Player Steam ID Struct returned by the steam API (userURL)
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

// Player represents a list of players BAN data
type Player struct {
	Players []playerData `json:"players"`
}

func updatePlayersIfNeeded(players map[string]Player, spreadsheetID string) {
	fmt.Println("M=updatePlayersIfNeeded Aperte Enter para continuar!")
	//reader := bufio.NewReader(os.Stdin)
	//reader.ReadString('\n')
	for idx, p := range players {
		data := p.Players[0]
		if data.VACBanned {
			fmt.Println("M=updatePlayersIfNeeded VAC " + data.SteamId)
			UpdateVACBanStatus(idx, data.DaysSinceLastBan, spreadsheetID)
		}
	}
}

// UnmarshalPlayerByID returns a player and its data obtained from Steam API
func UnmarshalPlayerByID(steamID string) Player {
	player := Player{}
	str, err := getPlayerStatus(steamID)
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

func getAllPlayersStatuses(userSteamID map[string]string) map[string]Player {
	fmt.Println("M=getAllPlayersStatuses Aperte Enter para continuar!")
	//reader := bufio.NewReader(os.Stdin)
	//reader.ReadString('\n')
	players := make(map[string]Player)
	for idx, value := range userSteamID {

		players[idx] = UnmarshalPlayerByID(value)
		fmt.Println(value)
	}
	return players
}

func getPlayerStatus(steamID string) ([]byte, error) {
	url := buildGetURL(steamID)
	log.Printf("M=getPlayerStatus url=%v\n", url)
	resp, err := http.Get(url)

	if err != nil {
		log.Printf("M=getPlayerStatus err=%s\n", err)
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
	log.Printf("M=getPlayerSteamID url=%v\n" + url)
	resp, err := http.Get(url)

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

func buildGetUserURL(userName string) string {
	fmt.Printf("M=buildGetUserURL userName=%s\n", userName)
	return userURL + valveKey + userParamKey + userName
}

func buildGetURL(steamID string) string {
	fmt.Printf("M=buildGetURL SteamID=%s\n", steamID)
	return vacBanURL + valveKey + steamIDParamKey + steamID
}

//UpdatePlayersStatus updates players status on fiver spreadsheet
func UpdatePlayersStatus(spreadsheetID string) {
	fmt.Printf("M=UpdatePlayersStatus spreadsheetID=%v Aperte Enter para continuar!\n", spreadsheetID)
	//reader := bufio.NewReader(os.Stdin)
	//reader.ReadString('\n')
	userSteamID := GetSteamIDs(spreadsheetID)

	players := getAllPlayersStatuses(userSteamID)

	updatePlayersIfNeeded(players, spreadsheetID)
}
