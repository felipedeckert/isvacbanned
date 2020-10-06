package service

import (
	"bytes"
	"io/ioutil"
	"isvacbanned/handler"
	"isvacbanned/messenger"
	"isvacbanned/mock"
	"isvacbanned/model"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	tb "gopkg.in/tucnak/telebot.v2"
)

func initMocks() {
	UserServiceClient = UserServiceMock{}
	UrlServiceClient = UrlServiceMock{}
	PlayerServiceClient = PlayerServiceMock{}
	handler.FollowHandlerClient = handler.FollowHandlerMock{}
	messenger.MessengerClient = messenger.MessengerMock{}
	model.FollowModelClient = model.FollowMock{}
	model.UserModelClient = model.UserMock{}
}

func TestSetUpFollowHandler(t *testing.T) {
	initMocks()

	user := tb.User{ID: 123, FirstName: "Gabriel", Username: "fallen"}

	chat := tb.Chat{ID: 456, FirstName: "Gabriel", Username: "fallen"}

	msg := tb.Message{Sender: &user, Payload: "12345678901234567", Chat: &chat}
	bot := tb.Bot{}

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

	expectedResult := int64(456)

	dbID := setUpFollowHandler(&msg, &bot)

	assert.EqualValues(t, expectedResult, dbID)
}
