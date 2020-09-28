package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeString(t *testing.T) {
	testStr := "test"
	resp := SanitizeString(testStr)

	assert.EqualValues(t, testStr, resp)

	testStr = "1 test"
	resp = SanitizeString(testStr)

	assert.EqualValues(t, testStr, resp)

	testStr = "🎈S e n s e 愛🖤"
	expectedResp := "S e n s e "
	resp = SanitizeString(testStr)

	assert.EqualValues(t, expectedResp, resp)

}
