package util

import (
	"strings"
	"unicode"
)

// Checks if a word or sentence is spelled the same backwards and forwards (with spaces removed).
// Returns FALSE if text is empty or not a palindrome.
func IsPalindrome(text string) bool {
	// Use runes to compare actual unicode values rather than strings.
	runes := []rune(strings.ToLower(stripSpaces(text)))
	if len := len(runes); len > 0 {
		// Loop through half of the of the list and compare the opposite
		// characters to see if they match.
		for i := 0; i < len/2; i++ {
			if runes[i] != runes[len-i-1] {
				return false
			}
		}
		return true
	}
	return false
}

// Strips out all whitespace-like characters from text.
// Lovingly borrowed from Stackoverflow
func stripSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}
