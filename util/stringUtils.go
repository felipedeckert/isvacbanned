package util

import (
	"log"
	"regexp"
)

// SanitizeString removes all non acsii chars form a string
func SanitizeString(input string) string {
	re, err := regexp.Compile(`[^\x00-\x7F]`)
	if err != nil {
		log.Fatal(err)
	}
	sanitizedInput := re.ReplaceAllString(input, "")

	if len(sanitizedInput) == 0 {
		sanitizedInput = "Unreadable nickname"
	}

	return sanitizedInput
}
