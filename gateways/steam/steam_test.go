package steam_test

import (
	"github.com/stretchr/testify/require"
	"isvacbanned/gateways/steam"
	"testing"
)

func TestGetArgumentFromURL(t *testing.T) {
	lastParam := "third"
	myURL := "https://first.com/second/" + lastParam
	res, err := steam.GetArgumentFromURL(myURL)

	require.NoError(t, err)
	require.Equal(t, lastParam, res)
}

func TestGetArgumentFromURLWithExtraSlash(t *testing.T) {
	lastParam := "third"
	myURL := "https://first.com/second/" + lastParam + "/"
	res, err := steam.GetArgumentFromURL(myURL)

	require.NoError(t, err)
	require.Equal(t, lastParam, res)
}

func TestGetArgumentFromURLInvalidURL(t *testing.T) {
	myURL := "httpfirst.comsecond"
	_, err := steam.GetArgumentFromURL(myURL)

	require.Error(t, err)
}
