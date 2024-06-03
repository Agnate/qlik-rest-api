package util

func IsPalindrome(text string) bool {
	runes := []rune(text)
	if len := len(runes); len > 0 {
		for i := 0; i < len/2; i++ {
			if runes[i] != runes[len-i-1] {
				return false
			}
		}
	}
	return true
}
