package moac_test

import (
	"errors"
	"math"
	"testing"

	"git.sr.ht/~seirdy/moac"
)

const margin = 0.0025 // acceptable error

type givensTestCase struct {
	expectedErrME error
	expectedErrBF error
	name          string
	given         moac.Givens
	expectedBF    float64
	expectedME    float64
	quantum       bool
}

func givensTestCases() []givensTestCase { //nolint:funlen // single statement; length only from test case count
	return []givensTestCase{
		{ // from README
			name:    "hitchhiker",
			quantum: true,
			given: moac.Givens{
				Mass:        5.97e24,
				Time:        1.45e17,
				Temperature: 1900,
				Password:    "ȣMǚHǨȎ#ŕģ=ʬƦQoţ}tʂŦȃťŇ+ħHǰĸȵʣɐɼŋĬŧǺʀǜǬɰ'ʮ0ʡěɱ6ȫŭ",
			},
			expectedBF:    1.401e-4,
			expectedME:    408.4,
			expectedErrBF: nil,
		},
		{ // same as above but without custom temp
			name:    "hitchhiker default temp",
			quantum: true,
			given: moac.Givens{
				Mass:     5.97e24,
				Time:     1.45e17,
				Password: "ȣMǚHǨȎ#ŕģ=ʬƦQoţ}tʂŦȃťŇ+ħHǰĸȵʣɐɼŋĬŧǺʀǜǬɰ'ʮ0ʡěɱ6ȫŭ",
			},
			expectedBF:    0.0986,
			expectedME:    427.3,
			expectedErrBF: nil,
		},
		{ // from blog post: https://seirdy.one/2021/01/12/password-strength.html
			name: "universe",
			given: moac.Givens{
				// default mass is the mass of the observable universe
				Entropy: 510,
			},
			expectedBF:    9.527e-62,
			expectedME:    307.3,
			expectedErrBF: nil,
		},
		{ // Should use the default provided entropy but fall back to the
			// lower computed value
			name: "only energy",
			given: moac.Givens{
				Energy: 4e52,
			},
			expectedBF:    0.0134,
			expectedME:    250,
			expectedErrBF: nil,
		},
		{
			name:          "Mising energy, mass",
			quantum:       false,
			given:         moac.Givens{},
			expectedBF:    0,
			expectedME:    307.3,
			expectedErrBF: moac.ErrMissingPE,
		},
		{
			name:    "Mising password",
			quantum: false,
			given: moac.Givens{
				Mass:             0,
				GuessesPerSecond: 0,
				Entropy:          0,
				Time:             0,
				Power:            0,
				EnergyPerGuess:   0,
			},
			expectedBF:    0,
			expectedME:    307.3,
			expectedErrBF: moac.ErrMissingPE,
		},
	}
}

func TestBruteForceability(t *testing.T) {
	for _, test := range givensTestCases() {
		t.Run(test.name, func(t *testing.T) {
			got, err := moac.BruteForceability(&test.given, test.quantum)
			if !errors.Is(err, test.expectedErrBF) {
				t.Fatalf("BruteForceability() = %v", err)
			}
			if beyondAcceptableMargin(got, test.expectedBF) {
				t.Errorf("Bruteforceability() = %.4g; want %.4g", got, test.expectedBF)
			}
		})
	}
}

func TestMinEntropy(t *testing.T) {
	for _, test := range givensTestCases() {
		t.Run(test.name, func(t *testing.T) {
			got, err := moac.MinEntropy(&test.given, test.quantum)

			if !errors.Is(err, test.expectedErrME) {
				t.Errorf("MinEntropy returned error %s, expected %s", err.Error(), test.expectedErrME.Error())
			}

			if beyondAcceptableMargin(got, test.expectedME) {
				t.Errorf("MinEntropy() = %.4g; want %.4g", got, test.expectedME)
			}
		})
	}
}

func beyondAcceptableMargin(got, expected float64) bool {
	return math.Abs(got-expected)/expected > margin
}
