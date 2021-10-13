package pwgen_test

import (
	"crypto/rand"
	"strings"
	"testing"

	"git.sr.ht/~seirdy/moac/v2/charsets"
	"git.sr.ht/~seirdy/moac/v2/pwgen"
)

func TestGenPwPanics(t *testing.T) {
	shouldPanic(t)
}

func shouldPanic(t *testing.T) {
	t.Helper()

	pwr := pwgen.PwRequirements{
		CharsetsWanted: charsets.ParseCharsets([]string{"ascii", "latin", "🦖؆ص😈"}),
		MinLen:         32,
	}

	csprng := rand.Reader
	rand.Reader = strings.NewReader("")

	defer func() {
		rand.Reader = csprng

		if out := recover(); out != "can't generate passwords: crypto/rand unavailable: EOF" {
			t.Errorf("panic due to CSPRNG unavailability sent unexpected message: %v", out)
		}
	}()

	_, _ = pwgen.GenPW(pwr) //nolint:errcheck // we're checking for panics; errors are checked elsewhere

	t.Errorf("pwgen should have panicked without access to a CSPRNG")
}
