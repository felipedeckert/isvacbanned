package service

import (
	"errors"
	"log"
	"strings"
)

func getArgumentFromURL(url string) (string, error) {
	splittedInput := strings.Split(url, "/")
	if len(splittedInput) < 2 {
		return "", errors.New("Invalid URL")
	}
	return splittedInput[len(splittedInput)-1], nil
}

func getSteamID(url string) (string, error) {
	steamID := url
	var err error
	var customID string
	if strings.Contains(url, "id") {
		customID, err = getArgumentFromURL(url)
		steamID = unmarshalPlayerByName(customID)
	} else if strings.Contains(url, "profile") {
		steamID, err = getArgumentFromURL(url)
	}

	if err != nil {
		log.Printf("M=getSteamID input=%v\n", url)

		return "", err
	}

	log.Printf("M=getSteamID input=%v argument=%v\n", url, steamID)

	return steamID, nil
}
