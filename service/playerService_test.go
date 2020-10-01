package service

import (
	"bytes"
	"io/ioutil"
	"isvacbanned/mock"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	Client = &mock.Client{}
}

func TestGetPlayerSteamID(t *testing.T) {
	steamID := "12345678901234567"

	// build response JSON
	myJSON := `{ "response": { "steamId":"12345678901234567",	"success": 1 } }`
	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(myJSON)))

	mock.GetFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	playerName := "fallen"
	res, err := getPlayerSteamID(playerName)

	assert.EqualValues(t, steamID, res)
	assert.Nil(t, err)
}

func TestGetPlayerCurrentNickname(t *testing.T) {
	steamID := "12345678901234567"
	expectedNickname := "fallen"
	// build response JSON
	myJSON := `{ "response" : { "players" : [ { "personaname" : "fallen" } ] } }`
	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(myJSON)))

	mock.GetFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	res := GetPlayerCurrentNickname(steamID)

	assert.EqualValues(t, expectedNickname, res)
}

func TestGetPlayerStatus(t *testing.T) {
	expectedSteamID := "12345678901234567"
	expectedVACBanStatus := false

	myJSON := `{ "players": [ { "SteamId": "12345678901234567", "CommunityBanned": false, "VACBanned": false, "NumberOfVACBans": 0, "DaysSinceLastBan": 0, "NumberOfGameBans": 0, "EconomyBan": "none" }	] }`

	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(myJSON)))

	mock.GetFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	res := GetPlayerStatus(expectedSteamID)

	assert.EqualValues(t, expectedSteamID, res.Players[0].SteamId)
	assert.EqualValues(t, expectedVACBanStatus, res.Players[0].VACBanned)
}

func TestGetAllPlayersStatus(t *testing.T) {
	expectedSteamID := "12345678901234567"
	expectedVACBanStatus := false

	myJSON := `{ "players": [ { "SteamId": "12345678901234567", "CommunityBanned": false, "VACBanned": false, "NumberOfVACBans": 0, "DaysSinceLastBan": 0, "NumberOfGameBans": 0, "EconomyBan": "none" } ] }`

	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(myJSON)))

	mock.GetFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	myMap := make(map[string]string)
	myMap["2"] = expectedSteamID

	res := getAllPlayersStatus(myMap)

	assert.EqualValues(t, expectedSteamID, res["2"].Players[0].SteamId)
	assert.EqualValues(t, expectedVACBanStatus, res["2"].Players[0].VACBanned)
}
