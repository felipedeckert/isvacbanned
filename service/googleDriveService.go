package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

//const spreadsheetID = "1yCglGeasnJtYsTrEvfGCoTM2syxmt_9_DwlsGeaITEc"
const yes = "YES"

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		log.Printf("M=tokenFromFile err=" + err.Error())
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

//GetSteamIDs return a map of row ID and steam ID of a giver spreadsheet
func GetSteamIDs(spreadsheetID string) map[string]string {
	fmt.Printf("M=GetSteamIDs step=START spreadsheetID=%v\n", spreadsheetID)

	srv := authenticate()

	readRange := "Cheaters!A2:D"
	fmt.Printf("M=GetSteamIDs spreadSheetID=%v\n", spreadsheetID)
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		fmt.Printf("Unable to retrieve data from sheet: %v\n", err)
	}

	if len(resp.Values) == 0 {
		log.Printf("No data found.")
		return nil
	}

	log.Printf("SteamID:")
	result := make(map[string]string)
	for _, row := range resp.Values {
		fmt.Printf("%v \n", row[0])
		fmt.Printf("%v \n", row[3])
		result[row[0].(string)] = row[3].(string)
	}
	return result
}

//UpdateVACBanStatus sets the vac ban status to YES, updates daysSinceLastBan and update date
func UpdateVACBanStatus(rowID string, daysSinceLastBan int, spreadsheetID string) {
	log.Printf("M=UpdateVACBanStatus step=START")
	srv := authenticate()

	var vr sheets.ValueRange
	dateFormat := "01/02/2006 15:04:05"
	currTime := time.Now().Format(dateFormat)
	myval := []interface{}{yes, strconv.Itoa(daysSinceLastBan), currTime}

	vr.Values = append(vr.Values, myval)

	writeRange := "J" + rowID + ":L"
	_, err := srv.Spreadsheets.Values.Update(spreadsheetID, writeRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		fmt.Printf("Unable to retrieve data from sheet: %v\n", err)
	}
	log.Printf("M=UpdateVACBanStatus step=END")
}

func authenticate() *sheets.Service {
	log.Printf("M=authenticate step=START")
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		fmt.Printf("Unable to read client secret file: %v\n", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		fmt.Printf("Unable to parse client secret file to config: %v\n", err)
	}
	client := getClient(config)

	srv, err := sheets.New(client)
	if err != nil {
		fmt.Printf("Unable to retrieve Sheets client: %v\n", err)
	}

	log.Printf("M=authenticate step=END")
	return srv
}
