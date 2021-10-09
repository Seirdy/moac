package pwgen_test

// Exhaustively test GenPW

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"testing"
	"unicode/utf8"

	"git.sr.ht/~seirdy/moac/v2/charsets"
	"git.sr.ht/~seirdy/moac/v2/entropy"
	"git.sr.ht/~seirdy/moac/v2/internal/bounds"
	"git.sr.ht/~seirdy/moac/v2/pwgen"
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

var ErrTooLong = fmt.Errorf("password too long: %w", pwgen.ErrInvalidLenBounds)

// Number of times to run each test-case.
// We run each test case multiple times because of the non-determinism inherent to GenPW().
const defaultLoops = 64

func buildTestCases(loops int) (testCases map[testGroupInfo][]pwgenTestCase, iterations int) {
	testCases, iterations = buildGoodTestCases(loops)
	testCases[testGroupInfo{name: "bad testcases"}] = buildBadTestCases()

	return testCases, iterations
}

func buildBadTestCases() []pwgenTestCase {
	return []pwgenTestCase{
		{
			name:           "too short for all charsets",
			charsetsWanted: []string{"lowercase", "uppercase", "numbers", "symbols", "latin", "🦖؆ص😈"},
			maxLen:         5,
			expectedErr:    pwgen.ErrInvalidLenBounds,
		},
		{
			name:           "too short for all ASCII",
			charsetsWanted: []string{"ascii"},
			maxLen:         3,
			expectedErr:    pwgen.ErrInvalidLenBounds,
		},
		{
			name:           "maxLen smaller than minLen",
			charsetsWanted: []string{"lowercase", "uppercase", "numbers", "symbols", "latin", "🦖؆ص😈"},
			maxLen:         12,
			minLen:         18,
			expectedErr:    pwgen.ErrInvalidLenBounds,
		},
		{
			name:           "maxLen is negative",
			charsetsWanted: []string{"lowercase", "uppercase", "numbers", "symbols", "latin", "🦖؆ص😈"},
			maxLen:         -18,
			minLen:         12,
			expectedErr:    bounds.ErrImpossibleNegative,
		},
		{
			name:           "minLen is negative",
			charsetsWanted: []string{"lowercase", "uppercase", "numbers", "symbols", "latin", "🦖؆ص😈"},
			maxLen:         18,
			minLen:         -12,
			expectedErr:    bounds.ErrImpossibleNegative,
		},
		{
			name:           "no characters",
			charsetsWanted: []string{""},
			maxLen:         18,
			minLen:         12,
			expectedErr:    pwgen.ErrInvalidLenBounds,
		},
	}
}

type testGroupInfo struct {
	name           string
	tooLongAllowed float64
}

type pwgenCharset struct {
	group          testGroupInfo
	charsetsWanted []string
}

func goodTestData() ([]pwgenCharset, []minMaxLen, []float64) {
	pwgenCharsets := []pwgenCharset{
		{
			group:          testGroupInfo{name: "everything"},
			charsetsWanted: []string{"lowercase", "uppercase", "numbers", "symbols", "latin", "世界🧛"},
		},
		{
			group:          testGroupInfo{name: "ascii"},
			charsetsWanted: []string{"ascii"},
		},
		{
			group:          testGroupInfo{name: "latin"},
			charsetsWanted: []string{"latin"},
		},
		{
			group: testGroupInfo{name: "nonprintable gibberish", tooLongAllowed: 23},
			charsetsWanted: []string{
				"uppercase", "numbers", "lowercase", `"O4UÞjÖÿ.ßòºÒ&Û¨5ü4äMî3îÌ`,
			},
		},
		{
			group: testGroupInfo{name: "tinyPassword"},
			charsetsWanted: []string{
				"uppercase", "numbers", "lowercase", "numbers", "numbers", "symbols", "lowercase", "ipaExtensions", "🧛",
			},
		},
		{
			group: testGroupInfo{name: "many dupe zeroes"},
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
			group: testGroupInfo{name: "complex custom charsets", tooLongAllowed: 0.3},
			charsetsWanted: []string{
				"uppercase", "numbers", "lowercase",
				"𓂸",
				"عظ؆ص",
				// lots of duplicate chars
				"ἀἁἂἃἄἅἆἇἈἉἊἋἌἍἎἏἐἑἒἓἔἕἘἙἚἛἜἝἠἡἢἣἤἥἦἧἨἩἪἫἬἭἮἯἰἱἲἳἴἵἶἷἸἹἺἻἼἽἾἿὀὁὂὃὄὅὈὉὊὋὌὍὐὑὒὓὔὕ" +
					"ὖὗὙὛὝὟὠὡὢὣὤὥὦὧὨὩὪὫὬὭὮὯὰάὲέὴήὶίὸόὺύὼώᾀᾁᾂᾃᾄᾅᾆᾇᾈᾉᾊᾋᾌᾍᾎᾏᾐᾑᾒᾓᾔᾕᾖᾗᾘᾙᾚᾛᾜᾝᾞᾟᾠᾡᾢᾣᾤᾥᾦᾧ" +
					"ᾨᾩᾪᾫᾬᾭᾮᾯᾰᾱᾲᾳᾴᾶᾷᾸᾹᾺΆᾼ᾽ι᾿῀῁ῂῃῄῆῇῈΈῊΉῌ῍῎῏ῐῑῒΐῖῗῘῙῚΊ῝῞῟ῠῡῢΰῤῥῦῧῨῩῪΎῬ῭΅`ῲῳῴῶῷῸΌῺΏῼ",
				"𓂸",
			},
		},
	}
	minMaxLengths := []minMaxLen{{0, 0}, {0, 32}, {0, 65537}, {80, 0}, {12, 50}, {0, 1}, {1, 1}, {12, 12}}
	entropiesWanted := []float64{0, 1, 32, 64, 256, 512}

	return pwgenCharsets, minMaxLengths, entropiesWanted
}

func calculateIteration(
	pwgenCharsets []pwgenCharset, minMaxLengths []minMaxLen, entropiesWanted []float64, loops int,
) (iterationsPerCharset int) {
	iterationsPerCharset = len(minMaxLengths) * len(entropiesWanted) * loops

	log.Printf(
		"running %d pwgen test cases %d times each, %d cases per pwgen charset in all.\n"+
			"each charset is run %d times in total",
		len(pwgenCharsets)*len(minMaxLengths)*len(entropiesWanted), loops,
		len(minMaxLengths)*len(entropiesWanted),
		iterationsPerCharset,
	)

	return iterationsPerCharset
}

func buildGoodTestCases(loops int) (testCases map[testGroupInfo][]pwgenTestCase, iterationsPerCharset int) {
	pwgenCharsets, minMaxLengths, entropiesWanted := goodTestData()
	iterationsPerCharset = calculateIteration(pwgenCharsets, minMaxLengths, entropiesWanted, loops)

	testCases = make(map[testGroupInfo][]pwgenTestCase, len(pwgenCharsets))

	var caseIndex int

	for _, entropyWanted := range entropiesWanted {
		for _, mml := range minMaxLengths {
			for i := range pwgenCharsets {
				charset := pwgenCharsets[i]
				newCase := buildTestCase(charset, mml, entropyWanted)

				testCases[charset.group] = append(
					testCases[charset.group],
					newCase,
				)
				caseIndex++
			}
		}
	}

	return testCases, iterationsPerCharset
}

func buildTestCase(charset pwgenCharset, mml minMaxLen, entropyWanted float64) pwgenTestCase {
	newCase := pwgenTestCase{
		expectedErr: nil, name: charset.group.name, charsetsWanted: charset.charsetsWanted,
		entropyWanted: entropyWanted, minLen: mml.minLen, maxLen: mml.maxLen,
	}

	if mml.maxLen > 0 && mml.maxLen < len(charsets.ParseCharsets(charset.charsetsWanted)) {
		newCase.expectedErr = pwgen.ErrInvalidLenBounds
	}

	return newCase
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

func pwUsesEachCharsetSinglePass(cs charsets.CharsetCollection, password []rune) (charsets.CharsetCollection, bool) {
	var (
		unusedCharsets charsets.CharsetCollection = make([]charsets.Charset, 0)
		pwCopy                                    = make([]rune, len(password))
		pass                                      = true
	)

	copy(pwCopy, password)

	for _, charset := range cs {
		pwCharIndex := latterUsesElemFromFormer(charset.Runes(), pwCopy)
		if pwCharIndex == -1 {
			unusedCharsets = append(unusedCharsets, charset)
			pass = false

			continue
		}

		pwCopy = append(pwCopy[:pwCharIndex], pwCopy[pwCharIndex+1:]...)
	}

	return unusedCharsets, pass
}

func pwUsesEachCharset(cs charsets.CharsetCollection, password []rune) error {
	if unusedCharsets, validPW := pwUsesEachCharsetSinglePass(cs, password); !validPW {
		if unusedCharset2, validPW2 := pwUsesEachCharsetSinglePass(unusedCharsets, password); !validPW2 {
			return errors.New(pwUsesEachCharsetErrStr(string(password), unusedCharset2))
		}
	}

	return nil
}

func pwUsesEachCharsetErrStr(password string, unusedCharsets charsets.CharsetCollection) string {
	var unusedCharsetsStr string

	for _, unusedCharset := range unusedCharsets {
		unusedCharsetsStr += unusedCharset.String()
		unusedCharsetsStr += "\n"
	}

	errorStr := fmt.Sprintf(
		"GenPW() = %s; didn't use %d charsets\nunused charsets: %s\n",
		password, len(unusedCharsets), unusedCharsetsStr,
	)

	return errorStr
}

func pwOnlyUsesCharsets(cs charsets.CharsetCollection, password []rune) (rune, bool) {
	allowedRunes := cs.Combined()

	for _, pwChar := range password {
		for i, allowedChar := range allowedRunes {
			if pwChar == allowedChar {
				break
			} else if i == len(allowedRunes)-1 {
				return pwChar, false
			}
		}
	}

	return ' ', true
}

func pwLongEnough(password string, minLen, maxLen int, entropyWanted float64) (float64, error) {
	entropyCalculated := entropy.Entropy(password)
	pwLen := utf8.RuneCountInString(password)

	if pwLen < minLen {
		return entropyCalculated, fmt.Errorf("generated pw length %d below min length %d", pwLen, minLen)
	}

	if entropyCalculated < entropyWanted {
		if pwLen < maxLen {
			return entropyCalculated, fmt.Errorf(
				"generated pw %s has entropy %.3g; wanted %.3g; password length below max",
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

func pwCorrectLength(pwRunes []rune, minLen, maxLen int, entropyWanted float64, cs charsets.CharsetCollection) error {
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
		_, truncatedUsesEachCharset := pwUsesEachCharsetSinglePass(cs, truncatedPass)

		truncatedEntropy := entropy.Entropy(string(truncatedPass))

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

func validateTestCase(test *pwgenTestCase, cs charsets.CharsetCollection) error {
	pwr := pwgen.PwRequirements{
		CharsetsWanted: charsets.ParseCharsets(test.charsetsWanted),
		TargetEntropy:  test.entropyWanted,
		MinLen:         test.minLen,
		MaxLen:         test.maxLen,
	}
	password, err := pwgen.GenPW(pwr)

	if unexpectedErr(err, test.expectedErr) {
		return fmt.Errorf("error in GenPW(): %w", err)
	}

	if err == nil && test.expectedErr != nil {
		return fmt.Errorf("expected error %w from GenPW, got nil", test.expectedErr)
	}

	pwRunes := []rune(password)
	if errUsesEachCharset := pwUsesEachCharset(cs, pwRunes); errUsesEachCharset != nil && err == nil {
		return errUsesEachCharset
	}

	if invalidRune, validPW := pwOnlyUsesCharsets(cs, pwRunes); !validPW {
		return fmt.Errorf("generated password %s used invalid character \"%v\"", password, string(invalidRune))
	}

	// skip this test if we expected the password to be generated successfully
	if err := pwCorrectLength(
		pwRunes, test.minLen, test.maxLen, test.entropyWanted, cs,
	); test.expectedErr == nil && err != nil {
		return fmt.Errorf("bad password length in test %v: %w", test.name, err)
	}

	return nil
}

func getLoops() int {
	loopsStr := os.Getenv("LOOPS")
	loops, err := strconv.ParseInt(loopsStr, 10, 64)

	if err != nil || loops == 0 {
		return defaultLoops
	}

	return int(loops)
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}

	return b
}

// need to split this into good/bad test cases.

func TestGenPw(t *testing.T) {
	t.Parallel()

	loops := getLoops()
	testCaseGroups, iterations := buildTestCases(loops)

	for name, testCaseGroup := range testCaseGroups {
		groupInfo, group := name, testCaseGroup

		allowedPercentWithOverage := groupInfo.tooLongAllowed
		if allowedPercentWithOverage > 0 {
			// with a low number of iterations, the percent overage is less
			// accurate so we need to be a bit more generous. The current
			// percent-overages are optimized for 100 loops and above, so become
			// more lenient while moving away from it.
			scaleFactor := 1 + math.Log(100/float64(loops))

			switch {
			case loops < 4:
				allowedPercentWithOverage = 100
			case loops < 10:
				allowedPercentWithOverage = (50 + groupInfo.tooLongAllowed) / 2
			case loops < 100:
				allowedPercentWithOverage = min(groupInfo.tooLongAllowed*scaleFactor, 33)
			default:
			}
		}

		t.Run(groupInfo.name, func(t *testing.T) {
			t.Parallel()

			tooLongCount := 0
			runTestCaseGroup(t, group, &tooLongCount, allowedPercentWithOverage, loops)
			log.Print(
				"number of too-long passwords for charset " +
					groupInfo.name +
					fmt.Sprintf(" %d/%d", tooLongCount, iterations),
			)

			if percent := float64(tooLongCount) / float64(iterations) * 100; percent > allowedPercentWithOverage {
				t.Errorf("%d out of %d passwords (%.3g%%) in charset group %s were too long; acceptable threshold is %.3g%%",
					tooLongCount, iterations, percent, groupInfo.name, allowedPercentWithOverage,
				)
			}
		})
	}
}

func TestGenPwHandlesSingleEmptyCharset(t *testing.T) {
	t.Parallel()

	pwr := pwgen.PwRequirements{
		CharsetsWanted: []charsets.Charset{charsets.CustomCharset(make([]rune, 0))},
		TargetEntropy:  128,
	}

	_, err := pwgen.GenPW(pwr)

	if !errors.Is(err, pwgen.ErrInvalidLenBounds) {
		t.Errorf("expected error %s from GenPW, got %s", pwgen.ErrInvalidLenBounds.Error(), err.Error())
	}
}

func runTestCaseGroup(
	t *testing.T, testCaseGroup []pwgenTestCase, tooLongCount *int, overageAllowed float64, loops int) {
	t.Helper()

	for i := range testCaseGroup {
		testCase := testCaseGroup[i]
		cs := charsets.ParseCharsets(testCase.charsetsWanted)

		for j := 0; j < loops; j++ {
			runTestCase(t, cs, &testCase, tooLongCount, overageAllowed)
		}
	}
}

func runTestCase(
	t *testing.T,
	cs charsets.CharsetCollection, testCase *pwgenTestCase, tooLongCount *int, overageAllowed float64,
) {
	t.Helper()

	err := validateTestCase(testCase, cs)
	if err != nil {
		if !errors.Is(err, ErrTooLong) || overageAllowed == 0 {
			t.Errorf(err.Error())
		}

		if *tooLongCount%3 == 0 && *tooLongCount < 15 { // don't spam output with >15 errors
			t.Log(err)
		}
		*tooLongCount++
	}
}
