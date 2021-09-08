// Package pwgen allows generating random passwords given charsets, length limits, and target entropy.
package pwgen

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	"strings"

	"git.sr.ht/~seirdy/moac"
	"git.sr.ht/~seirdy/moac/entropy"
)

func randInt(max int) int {
	newInt, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		log.Panicf("specialIndexes: %v", err)
	}

	return int(newInt.Int64())
}

func addRuneToPw(password *strings.Builder, runes []rune) {
	newChar := runes[randInt(len(runes))]
	password.WriteRune(newChar)
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

// computeSpecialIndexes determines the random locations at which to insert additional preselected chars.
// Generated passwords don't have truly uniform randomness; they also must have at
// least one of each charset, no matter how big/small that charset is. When we select
// one member of each charset, we need to insert those characters at random locations.
// specialIndexes determines those locations.
func computeSpecialIndexes(pwLength, charsetCount int) []int {
	res := make([]int, charsetCount)

	for i := 0; i < charsetCount; i++ {
		newInt := randInt(pwLength)

		for indexOf(res[0:i], newInt) >= 0 {
			newInt = randInt(pwLength)
		}

		res[i] = newInt
	}

	return res
}

func indexOf(src []int, e int) int {
	for i, a := range src {
		if a == e {
			return i
		}
	}

	return -1
}

func genpwFromGivenCharsets(charsetsGiven [][]rune, entropyWanted float64, minLen, maxLen int) (string, error) {
	var charsToPickFrom, pwBuilder strings.Builder

	if maxLen > 0 && maxLen < len(charsetsGiven) {
		return pwBuilder.String(), fmt.Errorf(
			"%w: maxLen too short to use all available charsets", ErrInvalidLenBounds,
		)
	}

	for _, charset := range charsetsGiven {
		charsToPickFrom.WriteString(string(charset))
	}

	runesToPickFrom := []rune(charsToPickFrom.String())
	// figure out the minimum acceptable length of the password
	// and fill that up before measuring entropy.
	pwLength, err := computePasswordLength(len(runesToPickFrom), entropyWanted, minLen, maxLen)
	if err != nil {
		return pwBuilder.String(), fmt.Errorf("can't generate password: %w", err)
	}

	if pwLength < len(charsetsGiven) {
		pwLength = len(charsetsGiven) // we know this is below maxLen
	}

	pwBuilder.Grow(pwLength + 1)

	specialIndexes := computeSpecialIndexes(pwLength, len(charsetsGiven))
	currentLength := 0

	for specialI := 0; currentLength < pwLength; currentLength++ {
		if i := indexOf(specialIndexes, currentLength); i >= 0 {
			addRuneToPw(&pwBuilder, charsetsGiven[i]) // one of each charset @ a special index
			specialI++
		} else {
			addRuneToPw(&pwBuilder, runesToPickFrom)
		}
	}

	pw := pwBuilder.String()
	pwRunes := []rune(pw)

	// keep inserting chars at random locations until the pw is long enough
	for ; maxLen == 0 || currentLength < maxLen; currentLength++ {
		newChar := runesToPickFrom[randInt(len(runesToPickFrom))]
		index := randInt(len(pwRunes))
		pwRunes = append(pwRunes[:index+1], pwRunes[index:]...)
		pwRunes[index] = newChar
		pw = string(pwRunes)

		computedEntropy, err := entropy.Entropy(pw)
		if err != nil {
			log.Panicf("failed to determine if password is long enough: %v", err)
		}

		if entropyWanted < computedEntropy {
			break
		}
	}

	return pw, nil
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
