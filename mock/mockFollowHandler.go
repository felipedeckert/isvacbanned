package mock

import tb "gopkg.in/tucnak/telebot.v2"

type FollowHandler struct {
	GetFollowHandler func(m *tb.Message, bot *tb.Bot, steamID, currNickname string, userID int64) int64
}

//FollowHandler is the mock client's `FollowHandler` func
func (f *FollowHandler) FollowHandler(m *tb.Message, bot *tb.Bot, steamID, currNickname string, userID int64) int64 {
	return GetFollowHandler(m, bot, steamID, currNickname, userID)
}

var (
	GetFollowHandler func(m *tb.Message, bot *tb.Bot, steamID, currNickname string, userID int64) int64
)
