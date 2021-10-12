package pwgen_test

import (
	"errors"
	"testing"

	"git.sr.ht/~seirdy/moac/v2/charsets"
	"git.sr.ht/~seirdy/moac/v2/pwgen"
)

// TestGenPwHandlesSuperLongPw runs few slow tests, so run it first.
func TestGenPwHandlesSuperLongPw(t *testing.T) {
	t.Parallel()

	pwr := pwgen.PwRequirements{
		CharsetsWanted: charsets.ParseCharsets([]string{"ascii", "latin", "🦖؆ص😈"}),
		MinLen:         getLoops() * getLoops(),
	}

	if _, err := pwgen.GenPW(pwr); err != nil {
		t.Errorf("error in GenPW: %s", err.Error())
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
