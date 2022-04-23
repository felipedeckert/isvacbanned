package entities

type PlayerInfo struct {
	Response PlayerNicknameData `json:"response"`
}

type PlayerNicknameData struct {
	Players []ResponseNicknameData `json:"players"`
}

type ResponseNicknameData struct {
	PersonaName string `json:"personaname"`
	SteamID string `json:"steamid"`
}
