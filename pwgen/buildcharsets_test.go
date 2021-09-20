package pwgen_test

import (
	"errors"
	"fmt"
	"testing"

	"git.sr.ht/~seirdy/moac/entropy"
	"git.sr.ht/~seirdy/moac/internal/slicing"
	"git.sr.ht/~seirdy/moac/pwgen"
)

type buildCharsetsTestCase struct {
	charsetsExpected map[string][]rune
	name             string
	charsetsNamed    []string
}

func buildCharsetsTables() []buildCharsetsTestCase { //nolint:funlen // single statement; length from tables
	return []buildCharsetsTestCase{
		{
			name: "many dupe numbers",
			charsetsNamed: []string{
				"uppercase", "lowercase", "numbers", "symbols", "latin", "🧛",
				"000000",
				"000000",
				"0",
				"000000000000000000000000000000",
				"1234000",
				"000000",
				"ascii",
				"000000000000000000000000000000",
				"898",
			},
			charsetsExpected: map[string][]rune{
				"symbols":        []rune("!\"#%&'()*+,-./:;<=>?@[\\]^_`{|}~$-"),
				"10":             []rune("1234"),
				"11":             []rune("0"),
				"lowercase":      []rune("abcdefghijklmnopqrstuvwxyz"),
				"latinExtendedA": []rune("ĀāĂăĄąĆćĈĉĊċČčĎďĐđĒēĔĕĖėĘęĚěĜĝĞğĠġĢģĤĥĦħĨĩĪīĬĭĮįİıĲĳĴĵĶķĸĹĺĻļĽľĿŀŁłŃńŅņŇňŉŊŋŌōŎŏŐőŒœŔŕŖŗŘřŚśŜŝŞşŠšŢţŤťŦŧŨũŪūŬŭŮůŰűŲųŴŵŶŷŸŹźŻżŽžſ"),
				"uppercase":      []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
				"numbers":        []rune("1234567"),
				"14":             []rune("89"),
				"latinExtendedB": []rune("ƀƁƂƃƄƅƆƇƈƉƊƋƌƍƎƏƐƑƒƓƔƕƖƗƘƙƚƛƜƝƞƟƠơƢƣƤƥƦƧƨƩƪƫƬƭƮƯưƱƲƳƴƵƶƷƸƹƺƻƼƽƾƿǀǁǂǃǄǅǆǇǈǉǊǋǌǍǎǏǐǑǒǓǔǕǖǗǘǙǚǛǜǝǞǟǠǡǢǣǤǥǦǧǨǩǪǫǬǭǮǯǰǱǲǳǴǵǶǷǸǹǺǻǼǽǾǿȀȁȂȃȄȅȆȇȈȉȊȋȌȍȎȏȐȑȒȓȔȕȖȗȘșȚțȜȝȞȟȠȡȢȣȤȥȦȧȨȩȪȫȬȭȮȯȰȱȲȳȴȵȶȷȸȹȺȻȼȽȾȿɀɁɂɃɄɅɆɇɈɉɊɋɌɍɎɏ"),
				"5":              []rune("🧛"),
				"ipaExtensions":  []rune("ɐɑɒɓɔɕɖɗɘəɚɛɜɝɞɟɠɡɢɣɤɥɦɧɨɩɪɫɬɭɮɯɰɱɲɳɴɵɶɷɸɹɺɻɼɽɾɿʀʁʂʃʄʅʆʇʈʉʊʋʌʍʎʏʐʑʒʓʔʕʖʗʘʙʚʛʜʝʞʟʠʡʢʣʤʥʦʧʨʩʪʫʬʭʮʯ"),
				"latin1":         []rune("¡¢£¤¥¦§¨©ª«¬®¯°±²³´μ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþÿ"),
			},
		},
		{
			name: "empty entries",
			charsetsNamed: []string{
				"uppercase", "lowercase", "numbers", "ascii",
				"000000",
				"000000",
				"",
				"000000000000000000000000000000",
				"",
				"000000",
				"0",
				"89",
			},
			charsetsExpected: map[string][]rune{
				"uppercase": []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
				"symbols":   []rune("!\"#%&'()*+,-./:;<=>?@[\\]^_`{|}~$-"),
				"4":         []rune("0"),
				"lowercase": []rune("abcdefghijklmnopqrstuvwxyz"),
				"numbers":   []rune("1234567"),
				"11":        []rune("89"),
			},
		},
		{
			name: "subset and composite",
			charsetsNamed: []string{
				string(entropy.Charsets["uppercase"]) + string(entropy.Charsets["lowercase"]),
				"lowercase", "numbers",
			},
			charsetsExpected: map[string][]rune{
				"lowercase": []rune("abcdefghijklmnopqrstuvwxyz"),
				"numbers":   []rune("0123456789"),
				"0":         []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"),
			},
		},
		{
			// say a user specifies a custom charset that shadows *almost* an entire previous charset, save for a single element
			// the only way to guarantee one from each set is to
			name: "composite missing one letter",
			charsetsNamed: []string{
				"lowercase", "numbers",
				"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxy",
			},
			charsetsExpected: map[string][]rune{
				"lowercase": entropy.Charsets["lowercase"],
				"numbers":   []rune("0123456789"),
				"0":         entropy.Charsets["uppercase"],
			},
		},
	}
}

var errBadCharset = errors.New("generated charsets don't match expected")

func formerNotFoundInLatter(former, latter map[string][]rune) [][]rune {
	formerSlice := slicing.MapToSlice(former)
	latterSlice := slicing.MapToSlice(latter)
	missing := [][]rune{}

	for _, v := range formerSlice {
		if !slicing.SliceContainsRuneSlice(latterSlice, v) {
			missing = append(missing, v)
		}
	}

	return missing
}

func expectedMatchesActual(expected, actual map[string][]rune) error {
	errStr := "%w:"

	errStrFirstHalf, pass := handleMissingCharsets(
		formerNotFoundInLatter(expected, actual), "missing expected entries")

	if !pass {
		errStr += "\n" + errStrFirstHalf
	}

	errStrSecondHalf, passSecondHalf := handleMissingCharsets(
		formerNotFoundInLatter(actual, expected), "contains unexpected entries")

	if !passSecondHalf {
		pass = false
		errStr += "\n" + errStrSecondHalf
	}

	if pass {
		return nil
	}

	return fmt.Errorf(errStr, errBadCharset)
}

func handleMissingCharsets(missing [][]rune, errType string) (errStr string, pass bool) {
	if len(missing) > 0 {
		errStr = "actual charset " + errType + ": ["
		for _, missingCharset := range missing {
			errStr += `"` + string(missingCharset) + `", `
		}

		errStr += "]"

		return errStr, false
	}

	return "", true
}

func TestBuildCharsets(t *testing.T) {
	for _, testCase := range buildCharsetsTables() {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// map order isn't deterministic, so repeat each test case a few times
			for i := 0; i < 1; i++ {
				charsetsActual := pwgen.BuildCharsets(testCase.charsetsNamed)
				if err := expectedMatchesActual(testCase.charsetsExpected, charsetsActual); err != nil {
					errStr := err.Error() + ":"
					for key, val := range charsetsActual {
						errStr += fmt.Sprintf("\n"+`"%s": []rune("%s"),`, key, string(val))
					}

					t.Error(errStr)

					break
				}
			}
		})
	}
}