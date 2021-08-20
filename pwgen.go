package moac

import (
	cryptoRand "crypto/rand"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"unicode/utf8"

	"git.sr.ht/~seirdy/moac/entropy"
)

func randRune(runes []rune) (rune, error) {
	i, err := cryptoRand.Int(cryptoRand.Reader, big.NewInt(int64(len(runes))))
	if err != nil {
		return ' ', fmt.Errorf("randRune: %w", err)
	}

	return runes[i.Int64()], nil
}

func addRuneToPw(password *string, runes []rune) error {
	newChar, err := randRune(runes)
	if err != nil {
		return fmt.Errorf("genpw: %w", err)
	}

	*password += string(newChar)

	return nil
}

func shuffle(password string) string {
	runified := []rune(password)
	rand.Shuffle(len(runified), func(i, j int) {
		runified[i], runified[j] = runified[j], runified[i]
	})

	return string(runified)
}

var errInvalidLenBounds = errors.New("bad length bounds")

func computePasswordLength(charsetSize int, pwEntropy float64, minLen, maxLen int) (int, error) {
	if maxLen > 0 && minLen > maxLen {
		return 0, fmt.Errorf("%w: maxLen can't be less than minLen", errInvalidLenBounds)
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
	var charsToPickFrom, pw string

	// at least one element from each charset
	for _, charset := range charsetsGiven {
		charsToPickFrom += string(charset)

		if err := addRuneToPw(&pw, charset); err != nil {
			return pw, fmt.Errorf("genpw: %w", err)
		}
	}

	runesToPickFrom := []rune(charsToPickFrom)
	// figure out the minimum acceptable length of the password and fill that up before measuring entropy.
	pwLength, err := computePasswordLength(len(runesToPickFrom), entropyWanted, minLen, maxLen)
	if err != nil {
		return "", fmt.Errorf("can't generate password: %w", err)
	}

	for utf8.RuneCountInString(pw) < pwLength {
		err := addRuneToPw(&pw, runesToPickFrom)
		if err != nil {
			return pw, fmt.Errorf("genpw: %w", err)
		}
	}

	for maxLen == 0 || utf8.RuneCountInString(pw) < maxLen {
		err := addRuneToPw(&pw, runesToPickFrom)
		if err != nil {
			return pw, fmt.Errorf("genpw: %w", err)
		}

		computedEntropy, err := entropy.Entropy(pw)
		if err != nil || entropyWanted < computedEntropy {
			return shuffle(pw), err
		}
	}

	return shuffle(pw), nil
}

func buildCharsets(charsetsEnumerated *[]string) [][]rune {
	var charsetsGiven [][]rune

	for _, charset := range *charsetsEnumerated {
		charsetRunes, found := entropy.Charsets[charset]

		switch {
		case found:
			charsetsGiven = append(charsetsGiven, charsetRunes)
		case charset == "latin":
			charsetsGiven = append(
				charsetsGiven,
				entropy.Charsets["latinExtendedA"], entropy.Charsets["latinExtendedB"], entropy.Charsets["ipaExtensions"],
			)
		default:
			charsetsGiven = append(charsetsGiven, []rune(charset))
		}
	}

	return charsetsGiven
}

// GenPW generates a random password using characters from the charsets enumerated by charsetsWanted.
// At least one element of each charset is used.
// Available charsets include "lowercase", "uppercase", "numbers", "symbols",
// "latinExtendedA", "latinExtendedB", and "ipaExtensions". "latin" is also available
// and is equivalent to specifying "latinExtendedA latinExtendedB ipaExtensions".
// Anything else will be treated as a string containing runes of a new custom charset
// to use.
// If entropyWanted is 0, the generated password has at least 256 bits of entropy;
// otherwise, it has entropyWanted bits of entropy.
// minLen and maxLen are ignored when set to zero; otherwise, they set lower/upper
// bounds on password character count and override entropyWanted if necessary.
func GenPW(charsetsEnumerated []string, entropyWanted float64, minLen, maxLen int) (string, error) {
	charsetsGiven := buildCharsets(&charsetsEnumerated)
	if entropyWanted == 0 {
		return genpwFromGivenCharsets(charsetsGiven, 256, minLen, maxLen)
	}

	return genpwFromGivenCharsets(charsetsGiven, entropyWanted, minLen, maxLen)
}
