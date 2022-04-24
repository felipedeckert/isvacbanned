package entities

type UsersFollowed struct {
	ID           int64
	SteamID      string
	OldNickname  string
	CurrNickname string
	IsCompleted  bool
}
