package moac_test

import (
	"math"
	"testing"

	"git.sr.ht/~seirdy/moac"
)

const margin = 0.05 // acceptable error

var givensTests = []struct {
	name    string
	given   moac.Givens
	quantum bool
	// Expected values should be within 10% error
	expectedBF float64
	expectedME float64
}{
	{ // from README
		name:    "hitchhiker",
		quantum: true,
		given: moac.Givens{
			Mass:     5.97e24,
			Time:     1.45e17,
			Password: "v¢JÊÙúQ§4mÀÛªZûYÍé©mËiÐ× \"½J6y.ñíí'è¦ïÏµ°~",
		},
		expectedBF: 0.6807,
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
}

func TestBruteForceability(t *testing.T) {
	for _, test := range givensTests {
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
	for _, test := range givensTests {
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
