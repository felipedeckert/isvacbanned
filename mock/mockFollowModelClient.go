package mock

import "isvacbanned/model"

//FollowModelClient is the mock client
type FollowModelClient struct {
	GetFollowSteamUser                func(chatID int64, steamID, currNickname string, userID int64) int64
	GetUnfollowSteamUser              func(userID int64, steamID string) int64
	GetGetFollowerCountBySteamID      func(steamID string) (int64, error)
	GetGetAllIncompletedFollowedUsers func() map[int64][]model.UsersFollowed
	GetGetUsersFollowed               func(userID int64) []model.UsersFollowed
	GetSetCurrNickname                func(userId int64, sanitizedActualNickname string)
	GetSetFollowedUserToCompleted     func(id []int64) int64
	GetIsFollowed                     func(steamID string, userID int64) (int64, error)
}

func (f *FollowModelClient) FollowSteamUser(chatID int64, steamID, currNickname string, userID int64) int64 {
	return GetFollowSteamUser(chatID, steamID, currNickname, userID)
}

func (f *FollowModelClient) UnfollowSteamUser(userID int64, steamID string) int64 {
	return GetUnfollowSteamUser(userID, steamID)
}

func (f *FollowModelClient) GetFollowerCountBySteamID(steamID string) (int64, error) {
	return GetGetFollowerCountBySteamID(steamID)
}

func (f *FollowModelClient) GetAllIncompletedFollowedUsers() map[int64][]model.UsersFollowed {
	return GetGetAllIncompletedFollowedUsers()
}

func (f *FollowModelClient) GetUsersFollowed(userID int64) []model.UsersFollowed {
	return GetGetUsersFollowed(userID)
}

func (f *FollowModelClient) SetCurrNickname(userId int64, sanitizedActualNickname string) {
	GetSetCurrNickname(userId, sanitizedActualNickname)
}

func (f *FollowModelClient) SetFollowedUserToCompleted(id []int64) int64 {
	return GetSetFollowedUserToCompleted(id)
}

func (f *FollowModelClient) IsFollowed(steamID string, userID int64) (int64, error) {
	return GetIsFollowed(steamID, userID)
}

var (
	GetFollowSteamUser                func(chatID int64, steamID, currNickname string, userID int64) int64
	GetUnfollowSteamUser              func(userID int64, steamID string) int64
	GetGetFollowerCountBySteamID      func(steamID string) (int64, error)
	GetGetAllIncompletedFollowedUsers func() map[int64][]model.UsersFollowed
	GetGetUsersFollowed               func(userID int64) []model.UsersFollowed
	GetSetCurrNickname                func(userId int64, sanitizedActualNickname string)
	GetSetFollowedUserToCompleted     func(id []int64) int64
	GetIsFollowed                     func(steamID string, userID int64) (int64, error)
)
