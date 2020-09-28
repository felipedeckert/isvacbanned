package service

import (
	"bytes"
	"io/ioutil"
	"isvacbanned/handler"
	"isvacbanned/mock"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	tb "gopkg.in/tucnak/telebot.v2"
)

func init() {
	userModelClient = &mock.UserModelClient{}
	handler.FollowClient = &mock.FollowModelClient{}
	handler.MsgClient = &mock.MsgClient{}
}

func TestSetUpFollowHandler(t *testing.T) {

	user := tb.User{ID: 123, FirstName: "Gabriel", Username: "fallen"}

	chat := tb.Chat{ID: 456}

	msg := tb.Message{Sender: &user, Payload: "12345678901234567", Chat: &chat}
	bot := tb.Bot{}

	mock.GetGetUserIDFunc = func(telegramID int) (int64, error) {
		return 123, nil
	}

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

	expectedResult := int64(321)

	mock.GetGetFollowerCountBySteamID = func(steamID string) (int64, error) {
		return 7, nil
	}

	mock.GetFollowSteamUser = func(chatID int64, steamID, currNickname string, userID int64) int64 {
		return expectedResult
	}

	mock.GetSendMessage = func(bot *tb.Bot, user *tb.User, message string) {}

	dbID := setUpFollowHandler(&msg, &bot)

	assert.EqualValues(t, expectedResult, dbID)
}
