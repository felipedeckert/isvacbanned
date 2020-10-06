package handler

import tb "gopkg.in/tucnak/telebot.v2"

type FollowHandlerMock struct{}

//HandleFollowRequest is the mock client's `HandleFollowRequest` func
func (f FollowHandlerMock) HandleFollowRequest(m *tb.Message, bot *tb.Bot, steamID, currNickname string, userID int64, isVACBanned bool) int64 {
	return int64(456)
}
