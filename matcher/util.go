package matcher

import (
	"log"
	"strings"
)

// startsWithAny checks whether any of words is the start of a sequence of words in the text
func anyWordStartsWith(text string, words ...string) bool {
	if len(words) == 0 {
		return false
	}

	var iterations = 0

	var currentIndex = 0

	for {
		iterations++

		var nextIndexOffset = strings.IndexFunc(text[currentIndex:], isAlphanumerical)
		if nextIndexOffset < 0 {
			break
		}
		currentIndex += nextIndexOffset

		for _, w := range words {
			if strings.HasPrefix(text[currentIndex:], w) {
				return true
			}
		}

		// Now skip to the next non-alphanumerical character
		nextIndexOffset = strings.IndexFunc(text[currentIndex:], func(r rune) bool {
			return !isAlphanumerical(r)
		})
		if nextIndexOffset < 0 {
			break
		}
		currentIndex += nextIndexOffset

		if iterations > 1000 {
			log.Printf("Input text %q causes containsAny to loop longer than expected", text)
			return false
		}
	}

	return false
}

func isAlphanumerical(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9')
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
