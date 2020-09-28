package handler

import "isvacbanned/model"

type FollowModelClient interface {
	FollowSteamUser(chatID int64, steamID, currNickname string, userID int64) int64
	GetFollowerCountBySteamID(steamID string) (int64, error)
	GetAllIncompletedFollowedUsers() map[int64][]model.UsersFollowed
	GetUsersFollowed(userID int64) []model.UsersFollowed
	SetCurrNickname(userId int, sanitizedActualNickname string)
	SetFollowedUserToCompleted(id []int) int64
}
