// Package entropy provides a means to compute entropy of a given random string
// by analyzing both the charsets used and its length.
package entropy

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"unicode/utf8"

	"git.sr.ht/~seirdy/moac/v2/charsets"
)

// Entropy computes the number of entropy bits in the given password,
// assumingly it was randomly generated.
func Entropy(password string) float64 {
	charsetsUsed := findCharsetsUsed(password)

	e, err := FromCharsets(charsetsUsed, utf8.RuneCountInString(password))
	// Should be impossible for FromCharsets to return an error when
	// charsetsUsed cannot be too long. If there's an error, we have a bug.
	if err != nil {
		log.Panicf("error measuring generated password entropy: %v", err)
	}

	return e
}

// FromCharsets computes the number of entropy bits in a string
// with the given length that utilizes at least one character from each
// of the given charsets. It does not perform any
// subsetting/de-duplication upon the given charsets; they are just used as-is.
func FromCharsets(charsetsUsed charsets.CharsetCollection, length int) (float64, error) {
	if len(charsetsUsed) > length {
		return 0.0, fmt.Errorf("password too short to use all available charsets: %w", ErrPasswordInvalid)
	}

	charSizeSum := 0

	for _, charset := range charsetsUsed {
		charSizeSum += len(charset.Runes())
	}
	// combos is charsize ^ length, entropy is ln2(combos)
	return float64(length) * math.Log2(float64(charSizeSum)), nil
}

func findCharsetsUsed(password string) charsets.CharsetCollection {
	var (
		filteredPassword = password
		charsetsUsed     charsets.CharsetCollection
	)

	for _, charset := range charsets.DefaultCharsets {
		if strings.ContainsAny(filteredPassword, charset.String()) {
			charsetsUsed.AddDefault(charset)
			filterFromString(&filteredPassword, charset.Runes())
		}
	}
	// any leftover characters that aren't from one of the hardcoded
	// charsets become a new charset of their own
	if filteredPassword != "" {
		charsetsUsed.Add(charsets.CustomCharset([]rune(filteredPassword)))
	}

	return charsetsUsed
}

func filterFromString(str *string, banned []rune) {
	*str = strings.Map(
		func(r rune) rune {
			for _, char := range banned {
				if char == r {
					return -1
				}
			}

			return r
		},
		*str,
	)
}

// A valid password is impossible with the given constraints.
var (
	ErrPasswordInvalid = errors.New("invalid password")
)
