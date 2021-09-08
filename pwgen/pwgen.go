// Package pwgen allows generating random passwords given charsets, length limits, and target entropy.
package pwgen

import (
	cryptoRand "crypto/rand"
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	"math/rand"
	"strings"
	"unicode/utf8"

	"git.sr.ht/~seirdy/moac"
	"git.sr.ht/~seirdy/moac/entropy"
)

func randRune(runes []rune) rune {
	i, err := cryptoRand.Int(cryptoRand.Reader, big.NewInt(int64(len(runes))))
	if err != nil {
		log.Panicf("crypto/rand errored when generating a random number: %v", err)
	}

	return runes[i.Int64()]
}

func addRuneToPw(password *strings.Builder, runes []rune) {
	newChar := randRune(runes)
	password.WriteRune(newChar)
}

func shuffle(password string) string {
	runified := []rune(password)
	rand.Shuffle(len(runified), func(i, j int) {
		runified[i], runified[j] = runified[j], runified[i]
	})

	return string(runified)
}

// ErrInvalidLenBounds represents bad minLen/maxLen values.
var ErrInvalidLenBounds = errors.New("bad length bounds")

func computePasswordLength(charsetSize int, pwEntropy float64, minLen, maxLen int) (int, error) {
	if maxLen > 0 && minLen > maxLen {
		return 0, fmt.Errorf("%w: maxLen can't be less than minLen", ErrInvalidLenBounds)
	}
	// combinations is 2^entropy, or 2^s
	// password length estimate is the logarithm of that with base charsetSize
	// logn(2^s) = s*logn(2) = s/log2(n)
	length := int(math.Ceil(pwEntropy / math.Log2(float64(charsetSize))))
	if length < minLen {
		length = minLen
	}

	if maxLen > 0 && length > maxLen {
		length = maxLen
	}

	return length, nil
}

func genpwFromGivenCharsets(charsetsGiven [][]rune, entropyWanted float64, minLen, maxLen int) (string, error) {
	var charsToPickFrom, pwBuilder strings.Builder

	// at least one element from each charset
	for _, charset := range charsetsGiven {
		charsToPickFrom.WriteString(string(charset))

		addRuneToPw(&pwBuilder, charset)
	}

	runesToPickFrom := []rune(charsToPickFrom.String())
	// figure out the minimum acceptable length of the password and fill that up before measuring entropy.
	pwLength, err := computePasswordLength(len(runesToPickFrom), entropyWanted, minLen, maxLen)
	if err != nil {
		return pwBuilder.String(), fmt.Errorf("can't generate password: %w", err)
	}

	pwBuilder.Grow(pwLength + 1)
	currentLength := utf8.RuneCountInString(pwBuilder.String())

	for ; currentLength < pwLength; currentLength++ {
		addRuneToPw(&pwBuilder, runesToPickFrom)
	}

	for ; maxLen == 0 || currentLength < maxLen; currentLength++ {
		addRuneToPw(&pwBuilder, runesToPickFrom)

		pw := pwBuilder.String()
		computedEntropy, err := entropy.Entropy(pw)

		if err != nil || entropyWanted < computedEntropy {
			return shuffle(pw), err
		}
	}

	return shuffle(pwBuilder.String()), nil
}

func buildCharsets(charsetsEnumerated []string) [][]rune {
	var charsetsGiven [][]rune

	for _, charset := range charsetsEnumerated {
		charsetRunes, found := entropy.Charsets[charset]

		switch {
		case found:
			charsetsGiven = append(charsetsGiven, charsetRunes)
		case charset == "latin":
			charsetsGiven = append(
				charsetsGiven,
				entropy.Charsets["latin1"],
				entropy.Charsets["latinExtendedA"],
				entropy.Charsets["latinExtendedB"],
				entropy.Charsets["ipaExtensions"],
			)
		default:
			charsetsGiven = append(charsetsGiven, []rune(charset))
		}
	}

	return charsetsGiven
}

// GenPW generates a random password using characters from the charsets enumerated by charsetsEnumerated.
// At least one element of each charset is used.
// Available charsets are "lowercase", "uppercase", "numbers", "symbols", "latin1",
// latinExtendedA", "latinExtendedB", and "ipaExtensions". "latin" is also available:
// it's equivalent to specifying "latin1 latinExtendedA latinExtendedB ipaExtensions".
// Anything else will be treated as a string containing runes of a new custom charset
// to use.
// If entropyWanted is 0, the generated password has at least 256 bits of entropy;
// otherwise, it has entropyWanted bits of entropy.
// minLen and maxLen are ignored when set to zero; otherwise, they set lower/upper
// bounds on password character count and override entropyWanted if necessary.
func GenPW(charsetsEnumerated []string, entropyWanted float64, minLen, maxLen int) (string, error) {
	charsetsGiven := buildCharsets(charsetsEnumerated)
	if entropyWanted == 0 {
		return genpwFromGivenCharsets(charsetsGiven, moac.DefaultEntropy, minLen, maxLen)
	}

	return genpwFromGivenCharsets(charsetsGiven, entropyWanted, minLen, maxLen)
}
