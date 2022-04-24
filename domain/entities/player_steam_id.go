package entities

// PlayerSteamID is the Player Steam ID Struct returned by the steam API (userURL)
type PlayerSteamID struct {
	Response responseData `json:"response"`
}

type responseData struct {
	SteamId string `json:"steamId"`
	Success int    `json:"success"`
}
