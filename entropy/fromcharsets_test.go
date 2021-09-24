package entropy_test

// This just tests that FromCharsets errors when the password length is
// too small to use all specified charsets.

import (
	"errors"
	"fmt"
	"testing"

	"git.sr.ht/~seirdy/moac/v2/charsets"
	"git.sr.ht/~seirdy/moac/v2/entropy"
)

type testCase struct {
	charsetsUsed charsets.CharsetCollection
	length       int
}

func buildTestCases() []testCase {
	return []testCase{
		{
			[]charsets.Charset{charsets.Lowercase},
			0,
		},
		{
			[]charsets.Charset{
				charsets.Lowercase, charsets.Uppercase,
				charsets.CustomCharset([]rune("¡¢£¤¥¦§¨©ª«¬®¯°±²³´μ¶·¸¹º»")),
			},
			2,
		},
		{
			[]charsets.Charset{
				charsets.Lowercase, charsets.Lowercase, charsets.Lowercase, charsets.Uppercase,
			},
			3,
		},
	}
}

func TestFromCharsetsErrors(t *testing.T) {
	for i, testCase := range buildTestCases() {
		testCase := testCase

		t.Run(fmt.Sprintf("fromCharsets case %d", i), func(t *testing.T) {
			_, err := entropy.FromCharsets(testCase.charsetsUsed, testCase.length)
			if !errors.Is(err, entropy.ErrPasswordInvalid) {
				t.Errorf("entropy.FromCharsets failed to error when given a password length that was too short")
			}
		})
	}
}
