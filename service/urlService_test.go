package service

import (
	"bytes"
	"io/ioutil"
	"isvacbanned/mock"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetArgumentFromURL(t *testing.T) {
	lastParam := "third"
	myURL := "http://first.com/second/" + lastParam
	res, err := getArgumentFromURL(myURL)

	if err != nil {
		t.Errorf("Unexpected erros err=%v", err.Error())
	}

	if res != lastParam {
		t.Errorf("got %v, want %v", res, lastParam)
	}
}

func TestGetArgumentFromURLWithExtraSlash(t *testing.T) {
	lastParam := "third"
	myURL := "http://first.com/second/" + lastParam + "/"
	res, err := getArgumentFromURL(myURL)

	if err != nil {
		t.Errorf("Unexpected erros err=%v", err.Error())
	}

	if res != lastParam {
		t.Errorf("got %v, want %v", res, lastParam)
	}
}

func TestGetArgumentFromURLInvalidURL(t *testing.T) {
	myURL := "httpfirst.comsecond"
	_, err := getArgumentFromURL(myURL)

	if err == nil {
		t.Errorf("Expected error, but no error returned")
	}
}

func TestGetSteamIDInvalidURL(t *testing.T) {
	lastParam := "third"
	myURL := "http://first.com/second/" + lastParam

	_, err := getSteamID(myURL)

	if err == nil {
		t.Errorf("Expected error, but no error returned")
	}
}

func TestGetSteamID(t *testing.T) {
	expectedParam := "third"
	myURL := "http://first.com/profile/" + expectedParam

	res, err := getSteamID(myURL)

	assert.Nil(t, err)

	assert.EqualValues(t, expectedParam, res)
}

func TestGetSteamIDFromCustomIDURL(t *testing.T) {
	expectedParam := "third"
	myURL := "http://first.com/id/" + expectedParam

	steamID := "12345678901234567"

	// build response JSON
	myJSON := `{ "response": { "steamId":"12345678901234567",	"success": 1 } }`
	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(myJSON)))

	mock.GetFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	res, err := getSteamID(myURL)

	assert.Nil(t, err)

	assert.EqualValues(t, steamID, res)
}
