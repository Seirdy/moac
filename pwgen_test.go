package moac // nolint:testpackage // use some private funcs cuz it's easier

import (
	"fmt"
	"log"
	"testing"
	"unicode/utf8"

	"git.sr.ht/~seirdy/moac/entropy"
)

type pwgenTestCase struct {
	name           string
	charsetsWanted []string
	entropyWanted  float64
	minLen         int
	maxLen         int
}

type minMaxLen struct {
	minLen int
	maxLen int
}

// Number of times to run each test-case.
// We run each test case multiple times because of the non-determinism inherent to GenPW().
const loops int = 16

var pwgenCharsets = []struct {
	name           string
	charsetsWanted []string
}{
	{
		name:           "everything",
		charsetsWanted: []string{"lowercase", "uppercase", "numbers", "symbols", "latin", "世界🧛"},
	},
	{
		name:           "alnum",
		charsetsWanted: []string{"lowercase", "uppercases", "numbers"},
	},
	{
		name: "tinyPassword",
		charsetsWanted: []string{
			"uppercase", "numbers", "lowercase", "numbers", "numbers", "symbols", "lowercase", "ipaExtensions", "🧛",
		},
	},
	{
		name: "multipleCustomCharsets",
		charsetsWanted: []string{
			"uppercase",
			"numbers",
			"lowercase",
			"𓂸",
			"عظ؆ص",
			"ἀἁἂἃἄἅἆἇἈἉἊἋἌἍἎἏἐἑἒἓἔἕἘἙἚἛἜἝἠἡἢἣἤἥἦἧἨἩἪἫἬἭἮἯἰἱἲἳἴἵἶἷἸἹἺἻἼἽἾἿὀὁὂὃὄὅὈὉὊὋὌὍὐὑὒὓὔὕὖὗὙὛὝὟὠὡὢὣὤὥὦὧὨὩὪὫὬὭὮὯὰάὲέὴήὶίὸόὺύὼώᾀᾁᾂᾃᾄᾅᾆᾇᾈᾉᾊᾋᾌᾍᾎᾏᾐᾑᾒᾓᾔᾕᾖᾗᾘᾙᾚᾛᾜᾝᾞᾟᾠᾡᾢᾣᾤᾥᾦᾧᾨᾩᾪᾫᾬᾭᾮᾯᾰᾱᾲᾳᾴᾶᾷᾸᾹᾺΆᾼ᾽ι᾿῀῁ῂῃῄῆῇῈΈῊΉῌ῍῎῏ῐῑῒΐῖῗῘῙῚΊ῝῞῟ῠῡῢΰῤῥῦῧῨῩῪΎῬ῭΅`ῲῳῴῶῷῸΌῺΏῼ", //nolint:lll
		},
	},
}

var (
	minMaxLengths   = []minMaxLen{{0, 0}, {0, 32}, {0, 65537}, {80, 0}, {12, 50}}
	entropiesWanted = []float64{0, 1, 32, 64, 256, 512}
)

func buildTestCases() []pwgenTestCase {
	log.Printf(
		"running %d pwgen test cases %d times each\n",
		len(pwgenCharsets)*len(minMaxLengths)*len(entropiesWanted), loops,
	)

	var (
		testCases []pwgenTestCase
		caseIndex int
	)

	for _, entropyWanted := range entropiesWanted {
		for _, minMaxLengths := range minMaxLengths {
			for _, pwgenCharset := range pwgenCharsets {
				testCases = append(
					testCases,
					pwgenTestCase{
						pwgenCharset.name, pwgenCharset.charsetsWanted,
						entropyWanted, minMaxLengths.minLen, minMaxLengths.maxLen,
					},
				)
				caseIndex++
			}
		}
	}

	return testCases
}

// second param should include at least one element of the first param.
func latterUsesFormer(former []rune, latter *[]rune) bool {
	for _, char := range former {
		for _, pwChar := range *latter {
			if pwChar == char {
				return true
			}
		}
	}

	return false
}

func pwUsesEachCharset(charsets *[][]rune, password *[]rune) (string, bool) {
	for _, charset := range *charsets {
		if !latterUsesFormer(charset, password) {
			return string(charset), false
		}
	}

	return "", true
}

func pwOnlyUsesAllowedRunes(charsets *[][]rune, password *[]rune) (rune, bool) {
	var allowedChars string
	for _, charset := range *charsets {
		allowedChars += string(charset)
	}

	allowedRunes := []rune(allowedChars)
	charSpace := len(allowedRunes)

	for _, pwChar := range *password {
		for i, char := range allowedRunes {
			if pwChar == char {
				break
			} else if i == charSpace-1 {
				return pwChar, false
			}
		}
	}

	return ' ', true
}

func pwHasGoodLength(password string, minLen, maxLen int, entropyWanted float64) error {
	entropyCalculated, err := entropy.Entropy(password)
	if err != nil {
		return fmt.Errorf("Error calculating entropy: %w", err)
	}

	pwLen := utf8.RuneCountInString(password)
	if maxLen > 0 && pwLen > maxLen {
		return fmt.Errorf("generated pw length %d exceeds max length %d", pwLen, maxLen)
	}

	if pwLen < minLen {
		return fmt.Errorf("generated pw length %d below min length %d", pwLen, minLen)
	}

	if entropyCalculated < entropyWanted {
		if pwLen < maxLen {
			return fmt.Errorf(
				"GenPW() = %s; entropy was %.3g, wanted %.3g, password length below max",
				password, entropyCalculated, entropyWanted,
			)
		}
	}

	return nil
}

func TestGenPw(t *testing.T) {
	for _, test := range buildTestCases() {
		charsetsWanted := &test.charsetsWanted
		entropyWanted := test.entropyWanted
		minLen := test.minLen
		maxLen := test.maxLen
		caseName := test.name
		t.Run(caseName, func(t *testing.T) {
			charsets := buildCharsets(charsetsWanted)
			for i := 0; i < loops; i++ {
				password, err := GenPW(*charsetsWanted, entropyWanted, minLen, maxLen)
				if err != nil {
					t.Fatalf("GenPW() = %v", err)
				}
				pwRunes := []rune(password)
				if unusedCharset, validPW := pwUsesEachCharset(&charsets, &pwRunes); !validPW {
					t.Errorf("GenPW() = %s; didn't use each charset\nunused charset: %s", password, unusedCharset)

					break
				}
				if invalidRune, validPW := pwOnlyUsesAllowedRunes(&charsets, &pwRunes); !validPW {
					t.Errorf("GenPW() = %s; used invalid character \"%v\"", password, string(invalidRune))

					break
				}
				if err != nil {
					t.Errorf("Error calculating entropy: %w", err)

					break
				}
				if err := pwHasGoodLength(password, minLen, maxLen, entropyWanted); err != nil {
					t.Errorf("bad password length in test %v: %w", caseName, err)

					break
				}
			}
		})
	}
}
