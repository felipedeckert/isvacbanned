package service

import (
	"errors"
	"log"
	"strings"
)

type urlService struct{}

type UrlServiceInterface interface {
	getSteamID(param string) (string, error)
}

var UrlServiceClient UrlServiceInterface = urlService{}

func (u urlService) getSteamID(param string) (string, error) {
	steamID := param
	var err error
	var customID string

	if strings.Contains(param, "id") { // URL with CustomID
		customID, err = getArgumentFromURL(param)
		if err != nil {
			log.Printf("M=getSteamID status=invalidCustomID param=%v\n", param)

			return "", err
		}
		steamID, err = getPlayerSteamID(customID)
		if err != nil {
			log.Printf("M=getSteamID status=CouldNotParseCustomID param=%v\n", param)

			return "", err
		}
	} else if strings.Contains(param, "profile") { // URL with SteamID
		steamID, err = getArgumentFromURL(param)
		if err != nil {
			log.Printf("M=getSteamID status=invalidSteamID param=%v\n", param)

			return "", err
		}
	} else { // CustomID without URL
		steamID, err = getPlayerSteamID(param)
		if err != nil {
			log.Printf("M=getSteamID status=notACustomID param=%v\n", param)
			return "", err
		}
	}

	log.Printf("M=getSteamID input=%v argument=%v\n", param, steamID)

	return steamID, nil
}

func getArgumentFromURL(url string) (string, error) {
	if last := len(url) - 1; last >= 0 && url[last] == '/' {
		url = url[:last]
	}

	splitInput := strings.Split(url, "/")
	if len(splitInput) < 2 {
		return "", errors.New("Invalid URL")
	}
	return splitInput[len(splitInput)-1], nil
}
