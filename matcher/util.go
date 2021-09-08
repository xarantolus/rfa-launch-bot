package matcher

import (
	"log"
	"strings"
	"unicode"
)

func anyWordStartsWith(text string, words ...string) bool {
	var iterations = 0

	var currentIndex = 0

	for {
		iterations++

		for currentIndex < len(text) && (unicode.IsSpace(rune(text[currentIndex])) || strings.ContainsAny(string(rune(text[currentIndex])), "#@$")) {
			currentIndex++
		}

		for _, w := range words {
			if strings.HasPrefix(text[currentIndex:], w) {
				return true
			}
		}

		// Now skip to the next space character
		for currentIndex < len(text) && !unicode.IsSpace(rune(text[currentIndex])) {
			currentIndex++
		}

		if currentIndex == len(text) {
			break
		}

		if iterations > 1000 {
			log.Printf("Input text %q causes containsAny to loop longer than expected", text)
			return false
		}
	}

	return false
}

// textContainsAny checks whether any of words is *anywhere* in the text
func textContainsAny(text string, words ...string) bool {
	for _, w := range words {
		if strings.Contains(text, w) {
			return true
		}
	}
	return false
}

// containsStringCaseInsensitive returns if the given slice contains word, comparing case-insensitive
func containsStringCaseInsensitive(slice []string, word string) bool {
	for _, w := range slice {
		if strings.EqualFold(w, word) {
			return true
		}
	}

	return false
}
