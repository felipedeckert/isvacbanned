package steam

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"isvacbanned/domain/entities"
	"isvacbanned/util"
	"log"
	"net/http"
	"strings"
)

//go:generate moq -stub -pkg mocks -out mocks/http_client.go . HTTPClient

const noMatch int = 42

type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}

type Steam struct {
	HTTPClient HTTPClient
}

func (s *Steam) GetSteamID(param string) (string, error) {
	steamID := param
	var err error
	var customID string

	if strings.Contains(param, "id") { // URL with CustomID
		customID, err = GetArgumentFromURL(param)
		if err != nil {
			log.Printf("M=getSteamID status=invalidCustomID param=%v\n", param)

			return "", err
		}
		steamID, err = s.GetPlayerSteamID(customID)
		if err != nil {
			log.Printf("M=getSteamID status=CouldNotParseCustomID param=%s\n", param)

			return "", err
		}
	} else if strings.Contains(param, "profile") { // URL with SteamID
		steamID, err = GetArgumentFromURL(param)
		if err != nil {
			log.Printf("M=getSteamID status=invalidSteamID param=%s\n", param)

			return "", err
		}
	} else {
		// CustomID without URL
		steamID, err = s.GetPlayerSteamID(param)
		if err != nil {
			log.Printf("M=getSteamID status=notACustomID param=%s\n", param)
			return "", err
		}
	}

	log.Printf("M=getSteamID input=%s argument=%s\n", param, steamID)

	return steamID, nil
}

func GetArgumentFromURL(url string) (string, error) {
	if last := len(url) - 1; last >= 0 && url[last] == '/' {
		url = url[:last]
	}

	splitInput := strings.Split(url, "/")
	if len(splitInput) < 2 {
		return "", errors.New("invalid URL")
	}
	return splitInput[len(splitInput)-1], nil
}

func (s *Steam) GetPlayerSteamID(playerName string) (string, error) {
	url := util.GetNicknameURL(playerName)
	log.Printf("M=getPlayerSteamID playerName=%v\n", playerName)

	resp, err := s.HTTPClient.Get(url)
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

	log.Printf("M=getPlayerSteamID steamID=%v\n", string(result))

	res, err := unmarshalSteamID(result)

	if err != nil {
		log.Printf("M=getPlayerSteamID step=unmarshal err=%s\n", err)
		return "", err
	}

	return res, nil
}

// GetPlayersCurrentNicknames gets the player identified by steamID current nickname
func (s *Steam) GetPlayersCurrentNicknames(steamIDs ...string) (entities.PlayerInfo, error) {
	url := util.GetPlayerSummaryURL(steamIDs...)
	resp, err := s.HTTPClient.Get(url)

	if err != nil {
		log.Printf("M=GetPlayersCurrentNickname err=%s\n", err)
		return entities.PlayerInfo{}, err
	}

	if resp.StatusCode != http.StatusOK {
		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf(`M=GetPlayersCurrentNickname L=I http status not OK, actual status:%d\n response:%s header:%v`, resp.StatusCode, string(result), resp.Header)
		return entities.PlayerInfo{}, err
	}

	result, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var playersInfo entities.PlayerInfo
	err = json.Unmarshal(result, &playersInfo)
	if err != nil {
		return entities.PlayerInfo{}, err
	}

	return playersInfo, nil
}

// GetPlayersStatus receives a list of steamIDs and returns the status for all players
func (s *Steam) GetPlayersStatus(steamIDs ...string) (entities.Player, error) {
	url := util.GetVACBanURL(steamIDs...)

	resp, err := s.HTTPClient.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf(`M=GetPlayerStatus L=I http status not OK, actual status:%d\n response:%s`, resp.StatusCode, string(result))
		return entities.Player{}, err
	}

	result, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	return unmarshalPlayer(result), nil
}

// unmarshalPlayer returns a player and its data obtained from steam API
func unmarshalPlayer(jsonInput []byte) entities.Player {
	player := entities.Player{}

	err := json.Unmarshal(jsonInput, &player)

	if err != nil {
		panic(err)
	}
	return player
}

func unmarshalSteamID(str []byte) (string, error) {
	playerID := entities.PlayerSteamID{}

	err := json.Unmarshal(str, &playerID)

	if err != nil {
		panic(err)
	}

	if playerID.Response.Success == noMatch {
		return "", errors.New("user not found")
	}

	return playerID.Response.SteamId, nil
}
