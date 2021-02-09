package handler

import tb "gopkg.in/tucnak/telebot.v2"

type UnfollowHandlerMock struct{}

//HandleFollowRequest is the mock client's `HandleFollowRequest` func
func (f UnfollowHandlerMock) HandleUnfollowRequest(m *tb.Message, bot *tb.Bot, steamID string, userID int64) {
}
