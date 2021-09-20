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

var ErrTooLong = fmt.Errorf("password too long: %w", ErrInvalidLenBounds)

// Number of times to run each test-case.
// We run each test case multiple times because of the non-determinism inherent to GenPW().
const loops int = 64

func buildTestCases() (testCases []pwgenTestCase, iterations int) {
	goodCases, iterations := buildGoodTestCases()

	return append(goodCases, buildBadTestCases()...), iterations
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
			name:           "too short for all ASCII",
			charsetsWanted: []string{"ascii"},
			maxLen:         3,
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
			name:           "ascii",
			charsetsWanted: []string{"ascii"},
		},
		{
			name:           "latin",
			charsetsWanted: []string{"latin"},
		},
		{
			name: "tinyPassword",
			charsetsWanted: []string{
				"uppercase", "numbers", "lowercase", "numbers", "numbers", "symbols", "lowercase", "ipaExtensions", "🧛",
			},
		},
		{
			name: "many dupe zeroes",
			charsetsWanted: []string{
				"uppercase", "lowercase", "numbers", "symbols", "latin", "🧛",
				"000000",
				"000000",
				"0",
				"000000000000000000000000000000",
				"1234000",
				"000000",
				"ascii",
				"000000000000000000000000000000",
			},
		},
		{
			name: "complex custom charsets",
			charsetsWanted: []string{
				"uppercase", "numbers", "lowercase",
				"𓂸",
				"عظ؆ص",
				// lots of duplicate chars
				"ἀἁἂἃἄἅἆἇἈἉἊἋἌἍἎἏἐἑἒἓἔἕἘἙἚἛἜἝἠἡἢἣἤἥἦἧἨἩἪἫἬἭἮἯἰἱἲἳἴἵἶἷἸἹἺἻἼἽἾἿὀὁὂὃὄὅὈὉὊὋὌὍὐὑὒὓὔὕὖὗὙὛὝὟὠὡὢὣὤὥὦὧὨὩὪὫὬὭὮὯὰάὲέὴήὶίὸόὺύὼώᾀᾁᾂᾃᾄᾅᾆᾇᾈᾉᾊᾋᾌᾍᾎᾏᾐᾑᾒᾓᾔᾕᾖᾗᾘᾙᾚᾛᾜᾝᾞᾟᾠᾡᾢᾣᾤᾥᾦᾧᾨᾩᾪᾫᾬᾭᾮᾯᾰᾱᾲᾳᾴᾶᾷᾸᾹᾺΆᾼ᾽ι᾿῀῁ῂῃῄῆῇῈΈῊΉῌ῍῎῏ῐῑῒΐῖῗῘῙῚΊ῝῞῟ῠῡῢΰῤῥῦῧῨῩῪΎῬ῭΅`ῲῳῴῶῷῸΌῺΏῼ", //nolint:lll // see prev
				"𓂸",
			},
		},
	}
	minMaxLengths := []minMaxLen{{0, 0}, {0, 32}, {0, 65537}, {80, 0}, {12, 50}, {0, 1}, {1, 1}, {12, 12}}
	entropiesWanted := []float64{0, 1, 32, 64, 256, 512}

	return pwgenCharsets, minMaxLengths, entropiesWanted
}

func buildGoodTestCases() (testCases []pwgenTestCase, iterationsPerCharset int) {
	pwgenCharsets, minMaxLengths, entropiesWanted := goodTestData()
	iterationsPerCharset = len(minMaxLengths) * len(entropiesWanted) * loops

	log.Printf(
		"running %d pwgen test cases %d times each, %d cases per pwgen charset in all.\n"+
			"each charset is run %d times in total",
		len(pwgenCharsets)*len(minMaxLengths)*len(entropiesWanted), loops,
		len(minMaxLengths)*len(entropiesWanted),
		iterationsPerCharset,
	)

	var caseIndex int

	for _, entropyWanted := range entropiesWanted {
		for _, minMaxLengths := range minMaxLengths {
			for _, pwgenCharset := range pwgenCharsets {
				newCase := pwgenTestCase{
					expectedErr: nil, name: pwgenCharset.name, charsetsWanted: pwgenCharset.charsetsWanted,
					entropyWanted: entropyWanted, minLen: minMaxLengths.minLen, maxLen: minMaxLengths.maxLen,
				}
				if minMaxLengths.maxLen > 0 && minMaxLengths.maxLen < len(BuildCharsets(pwgenCharset.charsetsWanted)) {
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

	return testCases, iterationsPerCharset
}

// second param should include at least one element of the first param.
func latterUsesElemFromFormer(former, latter []rune) int {
	for _, char := range former {
		for i, pwChar := range latter {
			if pwChar == char {
				return i
			}
		}
	}

	return -1
}

func pwUsesEachCharsetSinglePass(charsets map[string][]rune, password []rune) (map[string][]rune, bool) {
	var (
		unusedCharsets = make(map[string][]rune)
		pwCopy         = make([]rune, len(password))
		pass           = true
	)

	copy(pwCopy, password)

	for i := range charsets {
		pwCharIndex := latterUsesElemFromFormer(charsets[i], pwCopy)
		if pwCharIndex == -1 {
			unusedCharsets[i] = charsets[i]
			pass = false

			continue
		}

		pwCopy = append(pwCopy[:pwCharIndex], pwCopy[pwCharIndex+1:]...)
	}

	return unusedCharsets, pass
}

func pwUsesEachCharset(charsets map[string][]rune, password []rune) error {
	if unusedCharsets, validPW := pwUsesEachCharsetSinglePass(charsets, password); !validPW {
		if unusedCharset2, validPW2 := pwUsesEachCharsetSinglePass(unusedCharsets, password); !validPW2 {
			return errors.New(pwUsesEachCharsetErrStr(string(password), unusedCharset2))
		}
	}

	return nil
}

func pwUsesEachCharsetErrStr(password string, unusedCharsets map[string][]rune) string {
	var unusedCharsetsStr string

	for _, unusedCharset := range unusedCharsets {
		unusedCharsetsStr += string(unusedCharset)
		unusedCharsetsStr += "\n"
	}

	errorStr := fmt.Sprintf(
		"GenPW() = %s; didn't use %d charsets\nunused charsets: %s\n",
		password, len(unusedCharsets), unusedCharsetsStr,
	)

	return errorStr
}

func pwOnlyUsesAllowedRunes(charsets map[string][]rune, password *[]rune) (rune, bool) {
	var allowedChars string
	for _, charset := range charsets {
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

func pwLongEnough(password string, minLen, maxLen int, entropyWanted float64) (float64, error) {
	entropyCalculated, err := entropy.Entropy(password)
	if err != nil {
		return entropyCalculated, fmt.Errorf("error calculating entropy: %w", err)
	}

	pwLen := utf8.RuneCountInString(password)

	if pwLen < minLen {
		return entropyCalculated, fmt.Errorf("generated pw length %d below min length %d", pwLen, minLen)
	}

	if entropyCalculated < entropyWanted {
		if pwLen < maxLen {
			return entropyCalculated, fmt.Errorf(
				"GenPW() = %s; entropy was %.3g, wanted %.3g, password length below max",
				password, entropyCalculated, entropyWanted,
			)
		}
	}

	return entropyCalculated, nil
}

func unexpectedErr(actualErr, expectedErr error) bool {
	errorIsExpected := actualErr != nil && expectedErr != nil

	return errorIsExpected && !errors.Is(actualErr, expectedErr)
}

func pwCorrectLength(pwRunes []rune, minLen, maxLen int, entropyWanted float64, charsets map[string][]rune) error {
	pwLen := len(pwRunes)

	if maxLen > 0 && pwLen > maxLen {
		return fmt.Errorf("generated pw length %d exceeds max length %d", pwLen, maxLen)
	}

	entropyCalculated, err := pwLongEnough(string(pwRunes), minLen, maxLen, entropyWanted)
	if err != nil {
		return fmt.Errorf("failed to assert sufficient length: %w", err)
	}

	if pwLen > minLen && entropyWanted > 0 {
		truncatedPass := pwRunes[:len(pwRunes)-1]
		_, truncatedUsesEachCharset := pwUsesEachCharsetSinglePass(charsets, truncatedPass)

		truncatedEntropy, err := entropy.Entropy(string(truncatedPass))
		if err != nil {
			return fmt.Errorf("error calculating entropy: %w", err)
		}

		if truncatedEntropy >= entropyWanted && truncatedUsesEachCharset {
			return fmt.Errorf(
				"%w: "+
					"removing last char from password %v "+
					"caused its entropy to drop from %.4g to %.4g which is not below %.4g",
				ErrTooLong, string(pwRunes),
				entropyCalculated, truncatedEntropy, entropyWanted,
			)
		}
	}

	return nil
}

func validateTestCase(test *pwgenTestCase, charsets map[string][]rune) error {
	password, err := GenPW(test.charsetsWanted, test.entropyWanted, test.minLen, test.maxLen)
	if unexpectedErr(err, test.expectedErr) {
		return fmt.Errorf("GenPW() errored: %w", err)
	}

	if err == nil && test.expectedErr != nil {
		return fmt.Errorf("Expected error %w from GenPW, got nil", test.expectedErr)
	}

	pwRunes := []rune(password)
	if err == nil {
		if errUsesEachCharset := pwUsesEachCharset(charsets, pwRunes); errUsesEachCharset != nil {
			return errUsesEachCharset
		}
	}

	if invalidRune, validPW := pwOnlyUsesAllowedRunes(charsets, &pwRunes); !validPW {
		return fmt.Errorf("GenPW() = %s; used invalid character \"%v\"", password, string(invalidRune))
	}

	if unexpectedErr(err, test.expectedErr) {
		return fmt.Errorf("Error calculating entropy: %w", err)
	}

	// skip this test if we expected the password to be generated successfully
	if err := pwCorrectLength(
		pwRunes, test.minLen, test.maxLen, test.entropyWanted, charsets,
	); test.expectedErr == nil && err != nil {
		return fmt.Errorf("bad password length in test %v: %w", test.name, err)
	}

	return nil
}

func TestGenPw(t *testing.T) {
	var (
		testCases, iterations = buildTestCases()
		tooLongCount          int
	)

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			charsets := BuildCharsets(testCase.charsetsWanted)
			for j := 0; j < loops; j++ {
				err := validateTestCase(&testCase, charsets)
				if err != nil {
					if errors.Is(err, ErrTooLong) && testCase.name == "complex custom charsets" {
						log.Println(err)
						tooLongCount++
					} else {
						t.Error(err.Error())
					}

					break
				}
			}
		})
	}

	if percent := float64(tooLongCount) / float64(iterations) * 100; percent > 0.1 {
		t.Errorf("%d out of %d passwords (%.3g%%) in complex custom charsets were too long",
			tooLongCount, iterations, percent)
	}
}
