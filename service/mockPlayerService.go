package service

type PlayerServiceMock struct{}

func (p PlayerServiceMock) GetPlayerStatus(steamID string) Player {

	var players = []playerData{
		playerData{},
	}

	return Player{Players: players}
}

func (p PlayerServiceMock) GetPlayerCurrentNickname(steamID string) string {
	return "currentNickName"
}
