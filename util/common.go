package util

import (
	"log"
	"strings"
)

// LogError returns logs and returns an error, if err is nil it returns false
func LogError(err error, location string) bool {
	if err == nil {
		return false
	}
	log.Printf("[Error (%s)]: %s\n", location, err.Error())
	return true
}

func HashTagText(words []string) string {
	return "#" + strings.Join(words, " #")
}
