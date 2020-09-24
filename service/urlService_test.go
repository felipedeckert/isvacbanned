package service

import "testing"

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

func TestGetArgumentFromURLInvalidURL(t *testing.T) {
	myURL := "httpfirst.comsecond"
	_, err := getArgumentFromURL(myURL)

	if err == nil {
		t.Errorf("Expected error, but no error returned")
	}
}
