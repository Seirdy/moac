package entropy_test

// This just tests that FromCharsets errors when the password length is
// too small to use all specified charsets.

import (
	"errors"
	"fmt"
	"testing"

	"git.sr.ht/~seirdy/moac/entropy"
)

type testCase struct {
	charsetsUsed [][]rune
	length       int
}

func buildTestCases() []testCase {
	return []testCase{
		{
			[][]rune{[]rune("abcdefghijklmnopqrstuvwxyz")},
			0,
		},
		{
			[][]rune{entropy.Charsets["lowercase"], entropy.Charsets["uppercase"], []rune("¡¢£¤¥¦§¨©ª«¬®¯°±²³´μ¶·¸¹º»")},
			2,
		},
		{
			[][]rune{
				entropy.Charsets["lowercase"],
				entropy.Charsets["lowercase"],
				entropy.Charsets["lowercase"],
				entropy.Charsets["uppercase"],
			},
			3,
		},
	}
}

func TestFromCharsetsErrors(t *testing.T) {
	for i, testCase := range buildTestCases() {
		testCase := testCase

		t.Run(fmt.Sprintf("fromCharsets case %d", i), func(t *testing.T) {
			_, err := entropy.FromCharsets(&testCase.charsetsUsed, testCase.length)
			if !errors.Is(err, entropy.ErrPasswordInvalid) {
				t.Errorf("entropy.FromCharsets failed to error when given a password length that was too short")
			}
		})
	}
}
