package util

import "os"

const (
	vacBanURL = "http://api.steampowered.com/ISteamUser/GetPlayerBans/v1/?key="
	userURL = "http://api.steampowered.com/ISteamUser/ResolveVanityURL/v0001/?key="
 	playerSummaryURL = "http://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key="
 	userParamKey = "&vanityurl="
 	steamIDParamKey = "&steamids="
)

var valveKey string

func init() {
	valveKey = os.Getenv("STEAM_API_KEY")
}

//GetPlayerSummaryURL returns the URL to get the player's summary
func GetPlayerSummaryURL(steamID string) string {
	//todo change to debug level
	//log.Printf("M=buildGetPlayerSummaryURL SteamID=%s\n", steamID)
	return playerSummaryURL + valveKey + steamIDParamKey + steamID
}

//GetNicknameURL returns the URL to get the player's current nickname
func GetNicknameURL(userName string) string {
	//todo change to debug level
	//log.Printf("M=buildGetUserURL userName=%s\n", userName)
	return userURL + valveKey + userParamKey + userName
}

//GetVACBanURL returns the URL to get the player's ban status
func GetVACBanURL(steamID string) string {
	//todo change to debug level
	//log.Printf("M=buildGetURL steamID=%s\n", steamID)
	return vacBanURL + valveKey + steamIDParamKey + steamID
}
