package service

import (
	"database/sql"
	"testing"

	"isvacbanned/mock"

	"github.com/stretchr/testify/assert"
	tb "gopkg.in/tucnak/telebot.v2"
)

func init() {
	userModelClient = &mock.UserModelClient{}
}

func TestGetUserIDAlreadyCreated(t *testing.T) {

	userID := int64(123)

	user := tb.User{ID: 123, FirstName: "Gabriel", Username: "fallen"}

	mock.GetGetUserIDFunc = func(telegramID int) (int64, error) {
		return userID, nil
	}

	result := getUserID(&user)

	assert.EqualValues(t, userID, result)
}

func TestGetUserIDNewUser(t *testing.T) {

	userID := int64(123)

	user := tb.User{ID: 123, FirstName: "Gabriel", Username: "fallen"}

	mock.GetGetUserIDFunc = func(telegramID int) (int64, error) {
		return int64(-1), sql.ErrNoRows
	}

	mock.GetCreateUserFunc = func(firstName, username string, telegramID int) int64 {
		return userID
	}

	result := getUserID(&user)

	assert.EqualValues(t, userID, result)
}
