// Package sanitize implements detection and/or removal of unwanted characters in strings.
package sanitize

import (
	"strings"
	"unicode"
)

// FilterStrings removes bad runes from customCharsets, returning the result with the violating strings.
// sanitizedCharsets is customCharsets with all nonprintable or mark runes removed.
// badCharsets is the strings in customCharsets that contained such runes.
func FilterStrings(customCharsets []string) (sanitizedCharsets, badCharsets []string) {
	for _, charset := range customCharsets {
		sanitizedCharset, neededSanitization := filterCharset(charset)
		if neededSanitization {
			badCharsets = append(badCharsets, charset)
		}

		sanitizedCharsets = append(sanitizedCharsets, sanitizedCharset)
	}

	return sanitizedCharsets, badCharsets
}

func filterCharset(charset string) (string, bool) {
	var (
		filteredCharset strings.Builder
		containsBadRune bool
	)

	for _, r := range charset {
		if !goodRune(r) {
			containsBadRune = true

			continue
		}

		filteredCharset.WriteRune(r)
	}

	return filteredCharset.String(), containsBadRune
}

func goodRune(r rune) bool {
	return unicode.IsPrint(r) && !unicode.IsMark(r)
}
