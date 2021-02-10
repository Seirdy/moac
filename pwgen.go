package moac

import (
	cryptoRand "crypto/rand"
	"fmt"
	"math"
	"math/big"
	"math/rand"
)

const (
	lowercase     = "abcdefghijklmnopqrstuvwxyz"
	uppercase     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers       = "0123456789"
	symbols       = "!\"#%&'()*+,-./:;<=>?@[\\]^_`{|}~$-"
	extendedASCII = "¡¦§¨©«¬®¯°±´¶·¸»¿×÷¤¢£¥¹½¼²³¾ªáÁàÀâÂåÅäÄãÃæÆçÇðÐéÉèÈêÊëËíÍìÌîÎïÏñÑºóÓòÒôÔöÖõÕøØßúÚùÙûÛüÜýÝÿþÞµ"
)

func randRune(runes []rune) (rune, error) {
	i, err := cryptoRand.Int(cryptoRand.Reader, big.NewInt(int64(len(runes))))
	if err != nil {
		return ' ', fmt.Errorf("randRune: %w", err)
	}
	return runes[i.Int64()], nil
}

func shuffle(password string) string {
	runified := []rune(password)
	rand.Shuffle(len(runified), func(i, j int) {
		runified[i], runified[j] = runified[j], runified[i]
	})
	return string(runified)
}

// passwordLengthEstimate's results are slightly lower than the expected
// password length to allow pre-generating the first several characters
// of a password before slow entropy measurements
func passwordLengthEstimate(charsetSize int, entropy float64) int {
	// combinations is 2^entropy, or 2ⁿ
	// password length estimate is the logarithm of that with base charsetSize
	// logₛ(2ⁿ) = n*logₛ(2) = n/log₂(s)
	return int(entropy / math.Log2(float64(charsetSize)) * 0.8)
}

func genpwFromGivenCharsets(charsetsGiven [][]rune, entropy float64) (string, error) {
	var (
		charsToPickFrom string
		pw              string
	)
	// at least one element from each charset
	for _, charset := range charsetsGiven {
		charsToPickFrom += string(charset)
		newChar, err := randRune(charset)
		if err != nil {
			return pw, fmt.Errorf("genpw: %w", err)
		}
		pw += string(newChar)
	}
	runesToPickFrom := []rune(charsToPickFrom)
	minLength := passwordLengthEstimate(len(runesToPickFrom), entropy)
	for i := 0; i < minLength-len(charsetsGiven); i++ {
		newChar, err := randRune(runesToPickFrom)
		if err != nil {
			return pw, fmt.Errorf("genpw: %w", err)
		}
		pw += string(newChar)
	}
	pw = shuffle(pw)
	for i := 0; ; i++ {
		if calculateEntropy(pw) > entropy {
			break
		}
		newChar, err := randRune(runesToPickFrom)
		if err != nil {
			return pw, fmt.Errorf("genpw: %w", err)
		}
		pw += string(newChar)
		// shuffle every three character additions
		// so that we don't get one of each symbol crammed at the beginning;
		// that'd be less random.
		if i%3 == 0 {
			pw = shuffle(pw)
		}
	}
	return pw, nil
}

// GenPW generates a random password using characters from the charsets enumerated by charsetsWanted.
// At least one element of each charset is used.
// Available charsets include "lowercase", "uppercase", "numbers", "symbols", and "extendASCII".
// Anything else will be treated as a string containing elements of a new custom charset to use.
func GenPW(charsetsWanted []string, entropyWanted float64) (string, error) {
	var charsetsGiven [][]rune
	charsets := map[string][]rune{
		"lowercase":     []rune(lowercase),
		"uppercase":     []rune(uppercase),
		"numbers":       []rune(numbers),
		"symbols":       []rune(symbols),
		"extendedASCII": []rune(extendedASCII),
	}
	for _, charset := range charsetsWanted {
		if charsetRunes, found := charsets[charset]; found {
			charsetsGiven = append(charsetsGiven, charsetRunes)
		} else {
			charsetsGiven = append(charsetsGiven, []rune(charset))
		}
	}
	return genpwFromGivenCharsets(charsetsGiven, entropyWanted)
}
