package utils

import "strings"

// IncludesFuzzy now performs a case-insensitive substring match.
func IncludesFuzzy(str string, phrase string) bool {
	str = strings.ToLower(str)
	phrase = strings.ToLower(phrase)
	return strings.Contains(str, phrase)
}