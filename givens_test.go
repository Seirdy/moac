package moac_test

import (
	"errors"
	"math"
	"testing"

	"git.sr.ht/~seirdy/moac/v2"
	"git.sr.ht/~seirdy/moac/v2/internal/bounds"
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
			expectedBF:  7.288e-70,
			expectedBFQ: 1.467e-4,
			expectedME:  204.2,
		},
		{ // same as above but without custom temp
			name: "hitchhiker default temp",
			given: moac.Givens{
				Mass:     5.97e24,
				Time:     1.45e17,
				Password: "ȣMǚHǨȎ#ŕģ=ʬƦQoţ}tʂŦȃťŇ+ħHǰĸȵʣɐɼŋĬŧǺʀǜǬɰ'ʮ0ʡěɱ6ȫŭ",
			},
			expectedBF:  5.129e-67,
			expectedBFQ: 0.1032,
			expectedME:  213.7,
		},
		{
			name: "universe",
			given: moac.Givens{
				// default mass is the mass of the observable universe
				Entropy: 510,
			},
			expectedBF:  9.527e-62,
			expectedBFQ: 5.51e15,
			expectedME:  307.3,
		},
		{
			name: "solar dyson sphere", // test time/power being bottlenecks
			given: moac.Givens{
				Time:        1.45e17,
				Power:       3.828e26,
				Temperature: 1.5e7,
				Entropy:     198,
			},
			expectedBF:  0.962,
			expectedBFQ: 6.1e29,
			expectedME:  198,
		},
		{
			name: "solar dyson sphere 2: electric boogaloo",
			given: moac.Givens{
				Time:        1.45e17,
				Energy:      5.55e43, // now we use Givens.calculatePower()
				Temperature: 1.5e7,
				Entropy:     198,
			},
			expectedBF:  0.962,
			expectedBFQ: 6.1e29,
			expectedME:  198,
		},
		{
			name: "solar dyson sphere with GPS",
			given: moac.Givens{
				Time:             1.45e17,
				GuessesPerSecond: 2.67e42,
				Mass:             1.7e308,
				Temperature:      1.5e7,
				Entropy:          198,
			},
			expectedBF:  0.962,
			expectedBFQ: 6.1e29,
			expectedME:  198,
		},
		{
			name: "only energy",
			given: moac.Givens{
				Energy: 4.0e52,
			},
			expectedBF:  0.0134,
			expectedBFQ: 4.55e36,
			expectedME:  249.8,
		},
		{
			name: "impossibly high temp",
			given: moac.Givens{
				Energy:      4.0e52,
				Temperature: 1.5e32,
			},
			expectedErrBF: bounds.ErrImpossiblyHigh,
			expectedErrME: bounds.ErrImpossiblyHigh,
		},
		{
			name: "negativeTemp",
			given: moac.Givens{
				Energy:      4.0e52,
				Temperature: -1.0e-10,
			},
			expectedErrBF: bounds.ErrImpossibleNegative,
			expectedErrME: bounds.ErrImpossibleNegative,
		},
		{
			name: "negativeMass",
			given: moac.Givens{
				Energy: 4.0e52,
				Mass:   -1.0e-10,
			},
			expectedErrBF: bounds.ErrImpossibleNegative,
			expectedErrME: bounds.ErrImpossibleNegative,
		},
		{
			name: "negativePower",
			given: moac.Givens{
				Energy: 4.0e52,
				Power:  -1.0e-10,
			},
			expectedErrBF: bounds.ErrImpossibleNegative,
			expectedErrME: bounds.ErrImpossibleNegative,
		},
		{
			name: "negativeTime",
			given: moac.Givens{
				Energy: 4.0e52,
				Time:   -1.0e-10,
			},
			expectedErrBF: bounds.ErrImpossibleNegative,
			expectedErrME: bounds.ErrImpossibleNegative,
		},
		{
			name: "negativeGPS",
			given: moac.Givens{
				Energy:           4.0e52,
				GuessesPerSecond: -1.0e4,
			},
			expectedErrBF: bounds.ErrImpossibleNegative,
			expectedErrME: bounds.ErrImpossibleNegative,
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

	if err1 == nil && err2 == nil && expectedErr == nil {
		return
	}

	if !errors.Is(err1, expectedErr) {
		t.Errorf(
			`%s: got error "%v", expected "%v"`,
			funcName, err1, expectedErr,
		)
	}

	if !errors.Is(err2, expectedErr) {
		t.Errorf(
			`%s: got error "%v", expected "%v"`,
			funcName+"Quantum", err2, expectedErr,
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
	tcs := givensTestCases()
	for i := range tcs {
		t.Run(tcs[i].name, func(t *testing.T) {
			test := tcs[i]

			validateFunction(t, &test)
		})
	}
}

func TestMinEntropy(t *testing.T) {
	tcs := givensTestCases()
	for i := range tcs {
		t.Run(tcs[i].name, func(t *testing.T) {
			test := tcs[i]
			minEnt, errMinEnt := test.given.MinEntropy()
			minEntQ, errMinEntQ := test.given.MinEntropyQuantum()

			validateErrors(t, errMinEnt, errMinEntQ, test.expectedErrME, "MinEntropy")

			if beyondAcceptableMargin(minEnt, test.expectedME) {
				t.Errorf("MinEntropy() = %.4g; want %.4g", minEnt, test.expectedME)
			}

			if beyondAcceptableMargin(minEntQ, minEnt*2) {
				t.Errorf("MinEntropyQuantum() = %.4g; want %.4g", minEntQ, minEnt*2)
			}
		})
	}
}

func beyondAcceptableMargin(got, expected float64) bool {
	return math.Abs(got-expected)/expected > margin
}
