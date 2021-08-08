package moac

import (
	cryptoRand "crypto/rand"
	"fmt"
	"math"
	"math/big"
	"math/rand"

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

func computePasswordLength(charsetSize int, pwEntropy float64) int {
	// combinations is 2^entropy, or 2^s
	// password length estimate is the logarithm of that with base charsetSize
	// logn(2^s) = s*logn(2) = s/log2(n)
	return int(math.Ceil(pwEntropy / math.Log2(float64(charsetSize))))
}

func genpwFromGivenCharsets(charsetsGiven [][]rune, entropyWanted float64) (string, error) {
	var charsToPickFrom string
	pw := ""
	// at least one element from each charset
	for _, charset := range charsetsGiven {
		charsToPickFrom += string(charset)
		err := addRuneToPw(&pw, charset)
		if err != nil {
			return pw, fmt.Errorf("genpw: %w", err)
		}
	}
	runesToPickFrom := []rune(charsToPickFrom)
	// figure out the minimum length of the password and fill that up before measuring entropy.
	minLength := computePasswordLength(len(runesToPickFrom), entropyWanted)
	for i := 0; i < minLength-len(charsetsGiven); i++ {
		err := addRuneToPw(&pw, runesToPickFrom)
		if err != nil {
			return pw, fmt.Errorf("genpw: %w", err)
		}
	}
	for {
		err := addRuneToPw(&pw, runesToPickFrom)
		if err != nil {
			return pw, fmt.Errorf("genpw: %w", err)
		}
		computedEntropy, err := entropy.FromCharsets(&charsetsGiven, len(pw))
		if err != nil || entropyWanted < computedEntropy {
			return shuffle(pw), err
		}
	}
}

func buildCharsets(charsetsEnumerated *[]string) [][]rune {
	var charsetsGiven [][]rune
	for _, charset := range *charsetsEnumerated {
		if charsetRunes, found := entropy.Charsets[charset]; found {
			charsetsGiven = append(charsetsGiven, charsetRunes)
		} else if charset == "latin" {
			charsetsGiven = append(
				charsetsGiven,
				entropy.Charsets["latinExtendedA"], entropy.Charsets["latinExtendedB"], entropy.Charsets["ipaExtensions"],
			)
		} else {
			charsetsGiven = append(charsetsGiven, []rune(charset))
		}
	}
	return charsetsGiven
}

// GenPW generates a random password using characters from the charsets enumerated by charsetsWanted.
// At least one element of each charset is used.
// Available charsets include "lowercase", "uppercase", "numbers", "symbols", "latinExtendedA",
// "latinExtendedB", and "ipaExtensions". "latin" is also available and is equivalent to specifying
// "latinExtendedA latinExtendedB ipaExtensions". Anything else will be treated as a string
// containing runes of a new custom charset to use.
// If entropyWanted is 0, the generated password has at least 256 bits of entropy; otherwise, it
// has entropyWanted bits of entropy.
func GenPW(charsetsEnumerated []string, entropyWanted float64) (string, error) {
	charsetsGiven := buildCharsets(&charsetsEnumerated)
	if entropyWanted == 0 {
		return genpwFromGivenCharsets(charsetsGiven, 256)
	}
	return genpwFromGivenCharsets(charsetsGiven, entropyWanted)
}
