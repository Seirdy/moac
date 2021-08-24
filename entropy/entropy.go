// Package entropy provides a means to compute entropy of a given random string
// by analyzing both the charsets used and its length.
package entropy

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"unicode/utf8"
)

const (
	lowercase      = "abcdefghijklmnopqrstuvwxyz"
	uppercase      = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers        = "0123456789"
	symbols        = "!\"#%&'()*+,-./:;<=>?@[\\]^_`{|}~$-"
	latinExtendedA = "ĀāĂăĄąĆćĈĉĊċČčĎďĐđĒēĔĕĖėĘęĚěĜĝĞğĠġĢģĤĥĦħĨĩĪīĬĭĮįİıĲĳĴĵĶķĸĹĺĻļĽľĿŀŁłŃńŅņŇňŉŊŋŌōŎŏŐőŒœŔŕŖŗŘřŚśŜŝŞşŠšŢţŤťŦŧŨũŪūŬŭŮůŰűŲųŴŵŶŷŸŹźŻżŽžſ"                                                                                 //nolint:lll
	latinExtendedB = "ƀƁƂƃƄƅƆƇƈƉƊƋƌƍƎƏƐƑƒƓƔƕƖƗƘƙƚƛƜƝƞƟƠơƢƣƤƥƦƧƨƩƪƫƬƭƮƯưƱƲƳƴƵƶƷƸƹƺƻƼƽƾƿǀǁǂǃǄǅǆǇǈǉǊǋǌǍǎǏǐǑǒǓǔǕǖǗǘǙǚǛǜǝǞǟǠǡǢǣǤǥǦǧǨǩǪǫǬǭǮǯǰǱǲǳǴǵǶǷǸǹǺǻǼǽǾǿȀȁȂȃȄȅȆȇȈȉȊȋȌȍȎȏȐȑȒȓȔȕȖȗȘșȚțȜȝȞȟȠȡȢȣȤȥȦȧȨȩȪȫȬȭȮȯȰȱȲȳȴȵȶȷȸȹȺȻȼȽȾȿɀɁɂɃɄɅɆɇɈɉɊɋɌɍɎɏ" //nolint:lll
	ipaExtensions  = "ɐɑɒɓɔɕɖɗɘəɚɛɜɝɞɟɠɡɢɣɤɥɦɧɨɩɪɫɬɭɮɯɰɱɲɳɴɵɶɷɸɹɺɻɼɽɾɿʀʁʂʃʄʅʆʇʈʉʊʋʌʍʎʏʐʑʒʓʔʕʖʗʘʙʚʛʜʝʞʟʠʡʢʣʤʥʦʧʨʩʪʫʬʭʮʯ"                                                                                                                 //nolint:lll
)

// Charsets is a dictionary of known Unicode code blocks to use when generating passwords.
// All runes are printable and single-width.
var Charsets = map[string][]rune{ //nolint:gochecknoglobals // maps can't be const
	"lowercase":      []rune(lowercase),
	"uppercase":      []rune(uppercase),
	"numbers":        []rune(numbers),
	"symbols":        []rune(symbols),
	"latinExtendedA": []rune(latinExtendedA),
	"latinExtendedB": []rune(latinExtendedB),
	"ipaExtensions":  []rune(ipaExtensions),
}

// Entropy computes the number of entropy bits in the given password,
// assumingly it was randomly generated.
func Entropy(password string) (float64, error) {
	charsetsUsed := findCharsetsUsed(password)

	return FromCharsets(&charsetsUsed, utf8.RuneCountInString(password))
}

func findCharsetsUsed(password string) [][]rune {
	var (
		filteredPassword = password
		charsetsUsed     [][]rune
	)

	for _, charset := range Charsets {
		if strings.ContainsAny(filteredPassword, string(charset)) {
			charsetsUsed = append(charsetsUsed, charset)
			filterFromString(&filteredPassword, charset)
		}
	}
	// any leftover characters that aren't from one of the hardcoded
	// charsets become a new charset of their own
	if filteredPassword != "" {
		return append(charsetsUsed, []rune(filteredPassword))
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

var errPasswordInvalid = errors.New("invalid password")

// FromCharsets computes the number of entropy bits in a string
// with the given length that utilizes at least one character from each
// of the given charsets.
func FromCharsets(charsetsUsed *[][]rune, length int) (float64, error) {
	if len(*charsetsUsed) > length {
		return 0.0, fmt.Errorf("password too short to use all available charsets: %w", errPasswordInvalid)
	}

	charSizeSum := 0

	for _, charset := range *charsetsUsed {
		charSizeSum += len(charset)
	}
	// combos is charsize ^ length, entropy is ln2(combos)
	return float64(length) * math.Log2(float64(charSizeSum)), nil
}
