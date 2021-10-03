// Package pwgen allows generating random passwords given charsets, length limits, and target entropy.
package pwgen

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"git.sr.ht/~seirdy/moac/v2"
	"git.sr.ht/~seirdy/moac/v2/charsets"
	"git.sr.ht/~seirdy/moac/v2/entropy"
)

// PwRequirements holds the parameters used for password generation.
// These include charsets to pick from and target strength.
type PwRequirements struct {
	CharsetsWanted charsets.CharsetCollection
	TargetEntropy  float64
	MinLen, MaxLen int
}

// GenPW generates a random password using characters from the charsets enumerated by charsetsEnumerated.
// At least one element of each charset is used.
// If entropyWanted is 0, the generated password has at least 256 bits of entropy;
// otherwise, it has entropyWanted bits of entropy.
// minLen and maxLen are ignored when set to zero; otherwise, they set lower/upper
// bounds on password character count and override entropyWanted if necessary.
// GenPW will *not* strip any characters from given charsets that may be undesirable
// (newlines, control characters, etc.), and does not preserve grapheme clusters.
func GenPW(pwr PwRequirements) (string, error) {
	if pwr.MaxLen > 0 && pwr.MaxLen < len(pwr.CharsetsWanted) {
		return "", fmt.Errorf(
			"%w: maxLen too short to use all available charsets", ErrInvalidLenBounds,
		)
	}

	if pwr.TargetEntropy == 0 {
		return genpwFromGivenCharsets(PwRequirements{
			CharsetsWanted: pwr.CharsetsWanted,
			TargetEntropy:  moac.DefaultEntropy,
			MinLen:         pwr.MinLen, MaxLen: pwr.MaxLen,
		})
	}

	return genpwFromGivenCharsets(pwr)
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

		// must be unique
		for indexOf(res[0:i], newInt) >= 0 {
			newInt = randInt(pwLength)
		}

		res[i] = newInt
	}

	return res
}

func genpwFromGivenCharsets(pwr PwRequirements) (pw string, err error) {
	var (
		charsToPickFrom, pwBuilder strings.Builder
		pwUsesCustomCharset        bool
	)

	for _, charset := range pwr.CharsetsWanted {
		charsToPickFrom.WriteString(charset.String())

		if !pwUsesCustomCharset {
			pwUsesCustomCharset = charsets.IsDefault(charset)
		}
	}

	runesToPickFrom := []rune(charsToPickFrom.String())
	// figure out the minimum acceptable length of the password
	// and fill that up before measuring entropy.
	pwLength, err := computePasswordLength(len(runesToPickFrom), pwr.TargetEntropy, pwr.MinLen, pwr.MaxLen)
	if err != nil {
		return pwBuilder.String(), fmt.Errorf("can't generate password: %w", err)
	}

	if pwLength < len(pwr.CharsetsWanted) {
		pwLength = len(pwr.CharsetsWanted) // we know this is below maxLen
	}

	pwBuilder.Grow(pwLength + 1)

	specialIndexes := computeSpecialIndexes(pwLength, len(pwr.CharsetsWanted))
	pwRunes := buildFixedLengthPw(&pwBuilder, pwLength, specialIndexes, runesToPickFrom, pwr.CharsetsWanted)

	// keep inserting chars at random locations until the pw is long enough
	if pwUsesCustomCharset {
		randomlyInsertRunesTillStrong(pwr.MaxLen, &pwRunes, pwr.TargetEntropy, runesToPickFrom)
	}

	return string(pwRunes), nil
}

func buildFixedLengthPw(
	pwBuilder *strings.Builder,
	pwLength int, specialIndexes []int,
	runesToPickFrom []rune, cs charsets.CharsetCollection,
) []rune {
	currentLength := 0

	for specialI := 0; currentLength < pwLength; currentLength++ {
		if i := indexOf(specialIndexes, currentLength); i >= 0 {
			addRuneToEnd(pwBuilder, cs[i].Runes()) // one of each charset @ a special index
			specialI++
		} else {
			addRuneToEnd(pwBuilder, runesToPickFrom)
		}
	}

	return []rune(pwBuilder.String())
}

func randomlyInsertRunesTillStrong(maxLen int, pwRunes *[]rune, entropyWanted float64, runesToPickFrom []rune) {
	for maxLen == 0 || len(*pwRunes) < maxLen {
		if entropyWanted <= entropy.Entropy(string(*pwRunes)) {
			break
		}

		addRuneAtRandLoc(pwRunes, runesToPickFrom)
	}
}
