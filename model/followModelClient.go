package model

type FollowModelClient interface {
	FollowSteamUser(chatID int64, steamID, currNickname string, userID int64) int64
	UnfollowSteamUser(userID int64, steamID string) int64
	GetFollowerCountBySteamID(steamID string) (int64, error)
	GetAllIncompletedFollowedUsers() map[int64][]UsersFollowed
	GetUsersFollowed(userID int64) []UsersFollowed
	SetCurrNickname(userId int64, sanitizedActualNickname string)
	SetFollowedUserToCompleted(id []int64) int64
	IsFollowed(steamID string, userID int64) (int64, error)
}
