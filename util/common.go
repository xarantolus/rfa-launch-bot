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

func HashTagText(words []string) (s string) {
	var joined []string
	for _, w := range words {
		joined = append(joined, strings.Join(strings.Fields(w), ""))
	}
	return "#" + strings.Join(joined, " #")
}
