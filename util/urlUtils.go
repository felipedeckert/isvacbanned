package util

import "os"

const (
	vacBanURL        = "http://api.steampowered.com/ISteamUser/GetPlayerBans/v1/?key="
	userURL          = "http://api.steampowered.com/ISteamUser/ResolveVanityURL/v0001/?key="
	playerSummaryURL = "http://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key="
	userParamKey     = "&vanityurl="
	steamIDParamKey  = "&steamids="
	steamToken       = "STEAM_API_KEY"
)

//GetPlayerSummaryURL returns the URL to get the player's summary
func GetPlayerSummaryURL(steamIDs ...string) string {
	concatenatedIDs := ""

	for i, id := range steamIDs {
		concatenatedIDs += id
		if i + 1 < len(steamIDs) {
			concatenatedIDs += ","
		}
	}

	return playerSummaryURL + os.Getenv(steamToken) + steamIDParamKey + concatenatedIDs
}

//GetNicknameURL returns the URL to get the player's current nickname
func GetNicknameURL(userName string) string {

	return userURL + os.Getenv(steamToken) + userParamKey + userName
}

//GetVACBanURL returns the URL to get the player's ban status
func GetVACBanURL(steamIDs ...string) string {
	concatenatedIDs := ""

	for i, id := range steamIDs {
		concatenatedIDs += id
		if i + 1 < len(steamIDs) {
			concatenatedIDs += ","
		}
	}

	return vacBanURL + os.Getenv(steamToken) + steamIDParamKey + concatenatedIDs
}
