package model

type FollowMock struct{}

//FollowSteamUser is the mock implementation of FollowSteamUser
func (f FollowMock) FollowSteamUser(chatID int64, steamID, currNickname string, userID int64) int64 {
	return int64(99)
}

//UnfollowSteamUser is the mock implementation of UnfollowSteamUser
func (f FollowMock) UnfollowSteamUser(userID int64, steamID string) int64 {
	return int64(98)
}

//GetFollowerCountBySteamID is the mock implementation of GetFollowerCountBySteamID
func (f FollowMock) GetFollowerCountBySteamID(steamID string) (int64, error) {
	return int64(350), nil
}

//GetAllIncompletedFollowedUsers is the mock implementation of GetAllIncompletedFollowedUsers
func (f FollowMock) GetAllIncompletedFollowedUsers() map[int64][]UsersFollowed {
	return make(map[int64][]UsersFollowed, 0)
}

//GetUsersFollowed is the mock implementation of GetUsersFollowed
func (f FollowMock) GetUsersFollowed(userID int64) []UsersFollowed {
	return make([]UsersFollowed, 0)
}

//SetCurrNickname is the mock implementation of SetCurrNickname
func (f FollowMock) SetCurrNickname(userID int64, actualNickname string) {

}

//SetFollowedUserToCompleted is the mock implementation of SetFollowedUserToCompleted
func (f FollowMock) SetFollowedUserToCompleted(id []int64) int64 {
	return int64(97)
}

//GetUsersFollowedSummary is the mock implementation of GetUsersFollowedSummary
func (f FollowMock) GetUsersFollowedSummary(userID int64) map[bool]int {
	summary := make(map[bool]int, 0)

	summary[false] = 3
	summary[true] = 5
	return summary
}

//IsFollowed is the mock implementation of IsFollowed
func (f FollowMock) IsFollowed(steamID string, userID int64) (string, int64, error) {
	return "name", int64(1), nil
}
