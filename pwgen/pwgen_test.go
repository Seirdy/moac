package pwgen // nolint:testpackage // use some private funcs cuz it's easier

import (
	"errors"
	"fmt"
	"log"
	"testing"
	"unicode/utf8"

	"git.sr.ht/~seirdy/moac/entropy"
)

type pwgenTestCase struct {
	expectedErr    error
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
const loops int = 32

func buildTestCases() []pwgenTestCase {
	return append(buildGoodTestCases(), buildBadTestCases()...)
}

func buildBadTestCases() []pwgenTestCase {
	return []pwgenTestCase{
		{
			name:           "too short for all charsets",
			charsetsWanted: []string{"lowercase", "uppercase", "numbers", "symbols", "latin", "🦖؆ص😈"},
			maxLen:         5,
			expectedErr:    ErrInvalidLenBounds,
		},
		{
			name:           "bad lengths",
			charsetsWanted: []string{"lowercase", "uppercase", "numbers", "symbols", "latin", "🦖؆ص😈"},
			maxLen:         12,
			minLen:         18,
			expectedErr:    ErrInvalidLenBounds,
		},
	}
}

type pwgenCharset struct {
	name           string
	charsetsWanted []string
}

func goodTestData() ([]pwgenCharset, []minMaxLen, []float64) {
	pwgenCharsets := []pwgenCharset{
		{
			name:           "everything",
			charsetsWanted: []string{"lowercase", "uppercase", "numbers", "symbols", "latin", "世界🧛"},
		},
		{
			name:           "alnum",
			charsetsWanted: []string{"lowercase", "uppercase", "numbers"},
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
				"uppercase", "numbers", "lowercase",
				"𓂸",
				"عظ؆ص",
				// lots of duplicate chars
				"ἀἁἂἃἄἅἆἇἈἉἊἋἌἍἎἏἐἑἒἓἔἕἘἙἚἛἜἝἠἡἢἣἤἥἦἧἨἩἪἫἬἭἮἯἰἱἲἳἴἵἶἷἸἹἺἻἼἽἾἿὀὁὂὃὄὅὈὉὊὋὌὍὐὑὒὓὔὕὖὗὙὛὝὟὠὡὢὣὤὥὦὧὨὩὪὫὬὭὮὯὰάὲέὴήὶίὸόὺύὼώᾀᾁᾂᾃᾄᾅᾆᾇᾈᾉᾊᾋᾌᾍᾎᾏᾐᾑᾒᾓᾔᾕᾖᾗᾘᾙᾚᾛᾜᾝᾞᾟᾠᾡᾢᾣᾤᾥᾦᾧᾨᾩᾪᾫᾬᾭᾮᾯᾰᾱᾲᾳᾴᾶᾷᾸᾹᾺΆᾼ᾽ι᾿῀῁ῂῃῄῆῇῈΈῊΉῌ῍῎῏ῐῑῒΐῖῗῘῙῚΊ῝῞῟ῠῡῢΰῤῥῦῧῨῩῪΎῬ῭΅`ῲῳῴῶῷῸΌῺΏῼ", //nolint:lll // not worth splitting a single charset
				"𓂸",
			},
		},
	}
	minMaxLengths := []minMaxLen{{0, 0}, {0, 32}, {0, 65537}, {80, 0}, {12, 50}, {0, 1}, {1, 1}, {12, 12}}
	entropiesWanted := []float64{0, 1, 32, 64, 256, 512}

	return pwgenCharsets, minMaxLengths, entropiesWanted
}

func buildGoodTestCases() []pwgenTestCase {
	pwgenCharsets, minMaxLengths, entropiesWanted := goodTestData()

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
				newCase := pwgenTestCase{
					expectedErr: nil, name: pwgenCharset.name, charsetsWanted: pwgenCharset.charsetsWanted,
					entropyWanted: entropyWanted, minLen: minMaxLengths.minLen, maxLen: minMaxLengths.maxLen,
				}
				if minMaxLengths.maxLen > 0 && minMaxLengths.maxLen < len(pwgenCharset.charsetsWanted) {
					newCase.expectedErr = ErrInvalidLenBounds
				}

				testCases = append(
					testCases,
					newCase,
				)
				caseIndex++
			}
		}
	}

	return testCases
}

// second param should include at least one element of the first param.
func latterUsesFormer(former, latter []rune) bool {
	for _, char := range former {
		for _, pwChar := range latter {
			if pwChar == char {
				return true
			}
		}
	}

	return false
}

func pwUsesEachCharset(charsets [][]rune, password []rune) (string, bool) {
	for _, charset := range charsets {
		if !latterUsesFormer(charset, password) {
			return string(charset), false
		}
	}

	return "", true
}

func pwUsesEachCharsetErrStr(password, unusedCharset string, charsets [][]rune) string {
	errorStr := fmt.Sprintf(
		"GenPW() = %s; didn't use each charset\nunused charset: %s\ncharsets wanted are",
		password, unusedCharset,
	)
	for _, charset := range charsets {
		errorStr += "\n"
		errorStr += string(charset)
	}

	return errorStr
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
		return fmt.Errorf("error calculating entropy: %w", err)
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

func unexpectedErr(actualErr, expectedErr error) bool {
	errorIsExpected := actualErr != nil && expectedErr != nil

	return errorIsExpected && !errors.Is(actualErr, expectedErr)
}

func validateTestCase(test *pwgenTestCase, charsets [][]rune) error {
	password, err := GenPW(test.charsetsWanted, test.entropyWanted, test.minLen, test.maxLen)
	if unexpectedErr(err, test.expectedErr) {
		return fmt.Errorf("GenPW() errored: %w", err)
	}

	if err == nil && test.expectedErr != nil {
		return fmt.Errorf("Expected error %w from GenPW, got nil", test.expectedErr)
	}

	pwRunes := []rune(password)
	if err == nil {
		if unusedCharset, validPW := pwUsesEachCharset(charsets, pwRunes); !validPW {
			return fmt.Errorf(pwUsesEachCharsetErrStr(password, unusedCharset, charsets))
		}
	}

	if invalidRune, validPW := pwOnlyUsesAllowedRunes(&charsets, &pwRunes); !validPW {
		return fmt.Errorf("GenPW() = %s; used invalid character \"%v\"", password, string(invalidRune))
	}

	if unexpectedErr(err, test.expectedErr) {
		return fmt.Errorf("Error calculating entropy: %w", err)
	}

	// skip this test if we expected the password to be generated successfully
	if err := pwHasGoodLength(
		password, test.minLen, test.maxLen, test.entropyWanted,
	); test.expectedErr == nil && err != nil {
		return fmt.Errorf("bad password length in test %v: %w", test.name, err)
	}

	return nil
}

func TestGenPw(t *testing.T) {
	testCases := buildTestCases()
	for i := range testCases {
		t.Run(testCases[i].name, func(t *testing.T) {
			charsets := buildCharsets(testCases[i].charsetsWanted)
			for j := 0; j < loops; j++ {
				err := validateTestCase(&testCases[i], charsets)
				if err != nil {
					t.Error(err.Error())

					break
				}
			}
		})
	}
}
