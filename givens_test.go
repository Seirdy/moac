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
	expectedBFQ   float64
	expectedBF    float64
	expectedME    float64
}

// need to eventually remove the quantum bool from each test
// and instead test both quantum and non-quantum cases for each

func givensTestCases() []givensTestCase { //nolint:funlen // single statement; length only from test case count
	return []givensTestCase{
		{ // from README
			name: "hitchhiker",
			given: moac.Givens{
				Mass:        5.97e24,
				Time:        1.45e17,
				Temperature: 1900,
				Password:    "ȣMǚHǨȎ#ŕģ=ʬƦQoţ}tʂŦȃťŇ+ħHǰĸȵʣɐɼŋĬŧǺʀǜǬɰ'ʮ0ʡěɱ6ȫŭ",
			},
			expectedBF:  6.653e-70,
			expectedBFQ: 1.401e-4,
			expectedME:  204.2,
		},
		{ // same as above but without custom temp
			name: "hitchhiker default temp",
			given: moac.Givens{
				Mass:     5.97e24,
				Time:     1.45e17,
				Password: "ȣMǚHǨȎ#ŕģ=ʬƦQoţ}tʂŦȃťŇ+ħHǰĸȵʣɐɼŋĬŧǺʀǜǬɰ'ʮ0ʡěɱ6ȫŭ",
			},
			expectedBF:  4.682e-67,
			expectedBFQ: 0.0986,
			expectedME:  213.7,
		},
		{ //nolint:dupl // false positive
			name: "universe",
			given: moac.Givens{
				// default mass is the mass of the observable universe
				Entropy: 510,
			},
			expectedBF:  9.527e-62,
			expectedBFQ: 5.51e15,
			expectedME:  307.3,
		},
		{ //nolint:dupl // false positive
			name: "only energy",
			given: moac.Givens{
				Energy: 4e52,
			},
			expectedBF:  0.0134,
			expectedBFQ: 4.55e36,
			expectedME:  249.8,
		},
		{
			name:          "Mising energy, mass",
			given:         moac.Givens{},
			expectedBFQ:   0,
			expectedBF:    0,
			expectedME:    307.3,
			expectedErrBF: moac.ErrMissingPE,
		},
		{
			name: "Mising password",
			given: moac.Givens{
				Mass:             0,
				GuessesPerSecond: 0,
				Entropy:          0,
				Time:             0,
				Power:            0,
				EnergyPerGuess:   0,
			},
			expectedBFQ:   0,
			expectedBF:    0,
			expectedME:    307.3,
			expectedErrBF: moac.ErrMissingPE,
		},
	}
}

func validateErrors(t *testing.T, err1, err2, expectedErr error, funcName string) {
	t.Helper()

	if !errors.Is(errors.Unwrap(err1), errors.Unwrap(err2)) {
		t.Errorf(
			`%s: errors for non-quantum and quantum variants differ: "%s" != "%s"`,
			funcName, err1.Error(), err2.Error(),
		)
	}

	if !errors.Is(err1, expectedErr) {
		t.Errorf(
			`%s: got error "%s", expected "%s"`,
			funcName, err1.Error(), err2.Error(),
		)
	}
}

func validateFunction(t *testing.T, testCase *givensTestCase) {
	t.Helper()

	bf, errBF := testCase.given.BruteForceability()
	bfq, errBFQ := testCase.given.BruteForceabilityQuantum()

	validateErrors(t, errBF, errBFQ, testCase.expectedErrBF, "BruteForceability")

	if beyondAcceptableMargin(bf, testCase.expectedBF) {
		t.Errorf("BruteForceability() = %.4g; want %.4g", bf, testCase.expectedBF)
	}

	if beyondAcceptableMargin(bfq, testCase.expectedBFQ) {
		t.Errorf("BruteForceabilityQuantum() = %.4g; want %.4g", bfq, testCase.expectedBFQ)
	}
}

func TestBruteForceability(t *testing.T) {
	for _, test := range givensTestCases() {
		t.Run(test.name, func(t *testing.T) {
			test := test

			validateFunction(t, &test)
		})
	}
}

func TestMinEntropy(t *testing.T) {
	for _, test := range givensTestCases() {
		t.Run(test.name, func(t *testing.T) {
			me, errME := test.given.MinEntropy()
			meq, errMEQ := test.given.MinEntropyQuantum()

			validateErrors(t, errME, errMEQ, test.expectedErrME, "MinEntropy")

			if beyondAcceptableMargin(me, test.expectedME) {
				t.Errorf("MinEntropy() = %.4g; want %.4g", me, test.expectedME)
			}

			if beyondAcceptableMargin(meq, me*2) {
				t.Errorf("MinEntropyQuantum() = %.4g; want %.4g", meq, me*2)
			}
		})
	}
}

func beyondAcceptableMargin(got, expected float64) bool {
	return math.Abs(got-expected)/expected > margin
}
