package moac // nolint:testpackage // use some private funcs cuz it's easier

import (
	"fmt"
	"testing"
	"unicode/utf8"

	"git.sr.ht/~seirdy/moac/entropy"
)

var pwgenTests = []struct {
	name           string
	charsetsWanted []string
	entropyWanted  float64
}{
	{
		name:           "everything",
		charsetsWanted: []string{"lowercase", "uppercase", "numbers", "symbols", "latin", "世界🧛"},
		entropyWanted:  256,
	},
	{
		name:           "alnum",
		charsetsWanted: []string{"lowercase", "uppercases", "numbers"},
		entropyWanted:  64,
	},
	{
		name: "tinyPassword",
		charsetsWanted: []string{
			"uppercase", "numbers", "lowercase", "numbers", "numbers", "symbols", "lowercase", "ipaExtensions", "🧛",
		},
		entropyWanted: 1,
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
		entropyWanted: 256,
	},
}

type minMaxLen struct {
	minLen int
	maxLen int
}

var lengths = []minMaxLen{{0, 0}, {0, 32}, {0, 65537}}

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
			return fmt.Errorf("GenPW() = %s; entropy was %.3g, wanted %.3g, password length below max", password, entropyCalculated, entropyWanted)
		}
	}

	return nil
}

func TestGenPw(t *testing.T) {
	for _, test := range pwgenTests {
		charsetsWanted := &test.charsetsWanted
		entropyWanted := test.entropyWanted
		caseName := test.name
		t.Run(test.name, func(t *testing.T) {
			charsets := buildCharsets(charsetsWanted)
			// we're dealing with random passwords; try each testcase multiple times
			loops := 25
			for i := 0; i < loops; i++ {
				for _, pair := range lengths {
					password, err := GenPW(*charsetsWanted, entropyWanted, pair.minLen, pair.maxLen)
					if err != nil {
						t.Fatalf("GenPW() = %v", err)
					}
					pwRunes := []rune(password)
					if unusedCharset, validPW := pwUsesEachCharset(&charsets, &pwRunes); !validPW {
						t.Errorf("GenPW() = %s; didn't use each charset\nunused charset: %s", password, unusedCharset)
						i = loops
						break
					}
					if invalidRune, validPW := pwOnlyUsesAllowedRunes(&charsets, &pwRunes); !validPW {
						t.Errorf("GenPW() = %s; used invalid character \"%v\"", password, string(invalidRune))
						i = loops
						break
					}
					if err != nil {
						t.Errorf("Error calculating entropy: %w", err)
						i = loops
						break
					}
					if err := pwHasGoodLength(password, pair.minLen, pair.maxLen, entropyWanted); err != nil {
						t.Errorf("bad password length in test %v: %w", caseName, err)
						i = loops
						break
					}
				}
			}
		})
	}
}
