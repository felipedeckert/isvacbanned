package service

import (
	"testing"

	"isvacbanned/model"

	"github.com/stretchr/testify/assert"
	tb "gopkg.in/tucnak/telebot.v2"
)

func init() {
	model.UserModelClient = model.UserMock{}
}

func TestGetUserIDAlreadyCreated(t *testing.T) {

	userID := int64(321)

	user := tb.Chat{ID: 321, FirstName: "Gabriel", Username: "fallen"}

	result := UserServiceClient.getUserID(&user)

	assert.EqualValues(t, userID, result)
}

func TestGetUserIDNewUser(t *testing.T) {

	userID := int64(321)

	user := tb.Chat{ID: 321, FirstName: "Gabriel", Username: "fallen"}

	result := UserServiceClient.getUserID(&user)

	assert.EqualValues(t, userID, result)
}
