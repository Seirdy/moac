package moac_test

import (
	"math"
	"testing"

	"git.sr.ht/~seirdy/moac"
)

const margin = 0.0025 // acceptable error

type givensTestCase struct {
	name    string
	given   moac.Givens
	quantum bool
	// Expected values should be within 10% error
	expectedBF float64
	expectedME float64
}

func givensTestCases() []givensTestCase {
	return []givensTestCase{
		{ // from README
			name:    "hitchhiker",
			quantum: true,
			given: moac.Givens{
				Mass:     5.97e24,
				Time:     1.45e17,
				Password: "ȣMǚHǨȎ#ŕģ=ʬƦQoţ}tʂŦȃťŇ+ħHǰĸȵʣɐɼŋĬŧǺʀǜǬɰ'ʮ0ʡěɱ6ȫŭ",
			},
			expectedBF: 2.3e-49,
			expectedME: 427.3,
		},
		{ // from blog post: https://seirdy.one/2021/01/12/password-strength.html
			name:    "universe",
			quantum: false,
			given: moac.Givens{
				// default mass is the mass of the observable universe
				Entropy: 510,
			},
			expectedBF: 9.527e-62,
			expectedME: 307.3,
		},
		{ // Should use the default provided entropy but fall back to the
			// lower computed value
			name:    "only energy",
			quantum: false,
			given: moac.Givens{
				Energy: 4e52,
			},
			expectedBF: 0.0134,
			expectedME: 250,
		},
	}
}

func TestBruteForceability(t *testing.T) {
	for _, test := range givensTestCases() {
		t.Run(test.name, func(t *testing.T) {
			got, err := moac.BruteForceability(&test.given, test.quantum)
			if err != nil {
				t.Fatalf("BruteForceability() = %v", err)
			}
			if math.Abs(got-test.expectedBF)/test.expectedBF > margin {
				t.Errorf("Bruteforceability() = %.4g; want %.4g", got, test.expectedBF)
			}
		})
	}
}

func TestMinEntropy(t *testing.T) {
	for _, test := range givensTestCases() {
		t.Run(test.name, func(t *testing.T) {
			got, err := moac.MinEntropy(&test.given, test.quantum)
			if err != nil {
				t.Fatalf("MinEntropy() = %v", err)
			}
			if math.Abs(got-test.expectedME)/test.expectedME > margin {
				t.Errorf("MinEntropy() = %.4g; want %.4g", got, test.expectedME)
			}
		})
	}
}
