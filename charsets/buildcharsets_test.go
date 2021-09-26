package charsets_test

import (
	"testing"

	"git.sr.ht/~seirdy/moac/v2/charsets"
)

type buildCharsetsTestCase struct {
	charsetsExpected charsets.CharsetCollection
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
			charsetsExpected: stringsToCharsetCollection([]string{
				"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
				"abcdefghijklmnopqrstuvwxyz",
				"567",
				"!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~",
				"¡¢£¤¥¦§¨©ª«¬®¯°±²³´μ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþÿ",
				"ĀāĂăĄąĆćĈĉĊċČčĎďĐđĒēĔĕĖėĘęĚěĜĝĞğĠġĢģĤĥĦħĨĩĪīĬĭĮįİıĲĳĴĵĶķĸĹĺĻļĽľĿŀŁłŃńŅņŇňŉŊŋŌōŎŏŐőŒœŔŕŖŗŘřŚśŜŝŞşŠšŢţŤťŦŧŨũŪūŬŭŮůŰűŲųŴŵŶŷŸŹźŻżŽžſ",
				"ƀƁƂƃƄƅƆƇƈƉƊƋƌƍƎƏƐƑƒƓƔƕƖƗƘƙƚƛƜƝƞƟƠơƢƣƤƥƦƧƨƩƪƫƬƭƮƯưƱƲƳƴƵƶƷƸƹƺƻƼƽƾƿǀǁǂǃǄǅǆǇǈǉǊǋǌǍǎǏǐǑǒǓǔǕǖǗǘǙǚǛǜǝǞǟǠǡǢǣǤǥǦǧǨǩǪǫǬǭǮǯǰǱǲǳǴǵǶǷǸǹǺǻǼǽǾǿȀȁȂȃȄȅȆȇȈȉȊȋȌȍȎȏȐȑȒȓȔȕȖȗȘșȚțȜȝȞȟȠȡȢȣȤȥȦȧȨȩȪȫȬȭȮȯȰȱȲȳȴȵȶȷȸȹȺȻȼȽȾȿɀɁɂɃɄɅɆɇɈɉɊɋɌɍɎɏ",
				"ɐɑɒɓɔɕɖɗɘəɚɛɜɝɞɟɠɡɢɣɤɥɦɧɨɩɪɫɬɭɮɯɰɱɲɳɴɵɶɷɸɹɺɻɼɽɾɿʀʁʂʃʄʅʆʇʈʉʊʋʌʍʎʏʐʑʒʓʔʕʖʗʘʙʚʛʜʝʞʟʠʡʢʣʤʥʦʧʨʩʪʫʬʭʮʯ",
				"🧛",
				"0",
				"1234",
				"89",
			}),
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
			charsetsExpected: stringsToCharsetCollection([]string{
				"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
				"abcdefghijklmnopqrstuvwxyz",
				"1234567",
				"!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~",
				"0",
				"89",
			}),
		},
		{
			name: "unprintable gibberish",
			charsetsNamed: []string{
				charsets.Uppercase.String() + charsets.Lowercase.String(),
				"lowercase", "numbers", `"O4UÞjÖÿ.ßòºÒ&Û¨5ü4äMî3îÌ`,
			},
			charsetsExpected: stringsToCharsetCollection([]string{
				"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
				"abcdefghiklmnopqrstuvwxyz",
				"0123456789",
				`"&.j¨ºÌÒÖÛÞßäîòüÿ`,
			}),
		},
		{
			name: "subset and composite",
			charsetsNamed: []string{
				charsets.Uppercase.String() + charsets.Lowercase.String(),
				"lowercase", "numbers",
			},
			charsetsExpected: []charsets.Charset{charsets.Lowercase, charsets.Uppercase, charsets.Numbers},
		},
		{
			// say a user specifies a custom charset that shadows *almost* an entire previous charset, save for a single element
			// we shouldn't get a charset that includes only that single element
			name: "composite missing one letter",
			charsetsNamed: []string{
				"lowercase", "numbers",
				"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxy",
			},
			charsetsExpected: []charsets.Charset{charsets.Lowercase, charsets.Uppercase, charsets.Numbers},
		},
	}
}

func formerNotFoundInLatter(former, latter charsets.CharsetCollection) (missing charsets.CharsetCollection) {
	// for each charset in former: add to missing if it isn't in latter
	for _, f := range former {
		isMissing := true

		for _, l := range latter {
			if f.String() == l.String() {
				isMissing = false

				break
			}
		}

		if isMissing {
			missing = append(missing, f)
		}
	}

	return missing
}

func expectedMatchesActual(t *testing.T, expected, actual charsets.CharsetCollection) {
	t.Helper()

	var errStr string

	errStrFirstHalf, pass := handleMissingCharsets(
		formerNotFoundInLatter(expected, actual), "missing expected entries:   ")

	if !pass {
		errStr += "\n" + errStrFirstHalf
	}

	errStrSecondHalf, passSecondHalf := handleMissingCharsets(
		formerNotFoundInLatter(actual, expected), "contains unexpected entries:")

	if !passSecondHalf {
		pass = false
		errStr += "\n" + errStrSecondHalf
	}

	if pass {
		return
	}

	errStr = "generated charsets don't match expected:" + errStr + "\nactually got:"
	for _, actualCharset := range actual {
		errStr += "\n\"" + actualCharset.String() + `"`
	}

	if len(actual) == 0 {
		errStr = "built empty charset collection: " + errStr
	}

	t.Error(errStr)
}

func handleMissingCharsets(missing charsets.CharsetCollection, errType string) (errStr string, pass bool) {
	if len(missing) > 0 {
		errStr = "actual charset " + errType + " ["
		for _, missingCharset := range missing {
			errStr += `"` + missingCharset.String() + `", `
		}

		errStr += "]"

		return errStr, false
	}

	return "", true
}

func TestBuildCharsets(t *testing.T) {
	t.Parallel()

	for _, testCase := range buildCharsetsTables() {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			charsetsActual := charsets.ParseCharsets(testCase.charsetsNamed)
			expectedMatchesActual(t, testCase.charsetsExpected, charsetsActual)
		})
	}
}

func stringsToCharsetCollection(s []string) (cs charsets.CharsetCollection) {
	cs = make([]charsets.Charset, len(s))

	for i, ccContents := range s {
		cs[i] = charsets.CustomCharset([]rune(ccContents))
	}

	return cs
}
