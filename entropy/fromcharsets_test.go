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
	expectedErr  error
	charsetsUsed charsets.CharsetCollection
	length       int
}

func buildBadTestCases() []testCase {
	return []testCase{
		{
			charsetsUsed: []charsets.Charset{charsets.Lowercase},
			length:       0,
			expectedErr:  entropy.ErrPasswordInvalid,
		},
		{
			charsetsUsed: []charsets.Charset{
				charsets.Lowercase, charsets.Uppercase,
				charsets.CustomCharset([]rune("¡¢£¤¥¦§¨©ª«¬®¯°±²³´μ¶·¸¹º»")),
			},
			length:      2,
			expectedErr: entropy.ErrPasswordInvalid,
		},
		{
			charsetsUsed: []charsets.Charset{
				charsets.Lowercase, charsets.Lowercase, charsets.Lowercase, charsets.Uppercase,
			},
			length:      3, // FromCharsets does not perform any subsetting/duplication.
			expectedErr: entropy.ErrPasswordInvalid,
		},
	}
}

func buildGoodTestCases() []testCase {
	return []testCase{
		{
			charsetsUsed: []charsets.Charset{charsets.Lowercase},
			length:       1,
			expectedErr:  nil,
		},
		{
			charsetsUsed: []charsets.Charset{
				charsets.Lowercase, charsets.Uppercase,
				charsets.CustomCharset([]rune("¡¢£¤¥¦§¨©ª«¬®¯°±²³´μ¶·¸¹º»")),
			},
			length:      3,
			expectedErr: nil,
		},
		{
			charsetsUsed: []charsets.Charset{
				charsets.Lowercase, charsets.Lowercase, charsets.Lowercase, charsets.Uppercase,
			},
			length:      4,
			expectedErr: nil,
		},
	}
}

func TestBadCases(t *testing.T) {
	t.Parallel()

	tcs := append(buildBadTestCases(), buildGoodTestCases()...)
	for i, testCase := range tcs {
		testCase := testCase

		t.Run(fmt.Sprintf("fromCharsets case %d", i), func(t *testing.T) {
			t.Parallel()

			_, err := entropy.FromCharsets(testCase.charsetsUsed, testCase.length)
			validateTestCase(t, testCase.expectedErr, err)
		})
	}
}

func validateTestCase(t *testing.T, expectedErr, err error) {
	t.Helper()

	if expectedErr == nil {
		if err == nil {
			return
		}

		t.Errorf("entropy.FromCharsets returned unexpected error: %s", err.Error())
	}

	if !errors.Is(err, entropy.ErrPasswordInvalid) {
		t.Errorf("entropy.FromCharsets failed to error when given a password length that was too short")
	}
}
