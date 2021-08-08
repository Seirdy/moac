package moac

import (
	"errors"
	"fmt"
	"math"

	"github.com/nbutton23/zxcvbn-go"
)

// Givens holds the values used to compute password strength.
// These values are all physical quantities, measured using standard SI units.
type Givens struct {
	Password         string
	Entropy          float64
	Energy           float64
	Mass             float64 // mass used to build a computer or convert to energy
	Time             float64 // Duration of the attack, in seconds.
	EnergyPerGuess   float64
	Power            float64
	GuessesPerSecond float64
}

const (
	C         = 299792458      // speed of light in a vacuum, m/s
	G         = 6.67408e-11    // gravitation constant, m^3/kg/s^2
	Hubble    = 2.2e-18        // Hubble's Constant, hertz
	Temp      = 2.7            // cosmic radiation temperature (low estimate), kelvin
	Boltzmann = 1.3806503e-23  // Boltzmann's constant, J/K
	Planck    = 6.62607015e-35 // Planck's Constant, J*s

	UMass      = C * C * C / (2 * G * Hubble) // mass of the observable universe.
	Bremermann = C * C / Planck               // Bremermann's limit
	Landauer   = Boltzmann * Temp * math.Ln2  // Landauer limit

	defaultEntropy = 256
)

// populateDefaults fills in default values for entropy calculation if not provided
func populateDefaults(givens *Givens) {
	if givens.Entropy == 0 {
		if givens.Mass+givens.EnergyPerGuess == 0 {
			givens.Entropy = defaultEntropy
		} else if givens.Mass == 0 {
			// mass of the observable universe
			givens.Mass = UMass
		}
	}
	if givens.EnergyPerGuess == 0 {
		// maybe put something more elaborate here given different constraints
		givens.EnergyPerGuess = Landauer
	}

}

func calculateEntropy(password string) float64 {
	// currently wraps zxcvbn-go. This might change in the future.
	return zxcvbn.PasswordStrength(password, nil).Entropy
}

func calculatePower(givens *Givens) {
	powerFromComputationSpeed := givens.GuessesPerSecond * givens.EnergyPerGuess
	powerFromEnergy := givens.Energy / givens.Time
	// loop over an array for this since its length will grow in the future
	computedPowers := [2]float64{powerFromComputationSpeed, powerFromEnergy}
	for _, power := range computedPowers {
		if givens.Power == 0 || (power > 0 && power < givens.Power) {
			givens.Power = power
		}
	}
}

func calculateEnergy(givens *Givens) {
	massEnergy := givens.Mass * C * C
	energyFromPower := givens.Power * givens.Time
	computedEnergies := [2]float64{massEnergy, energyFromPower}
	for _, energy := range computedEnergies {
		if givens.Energy == 0 || (energy > 0 && energy < givens.Energy) {
			givens.Energy = energy
		}
	}
}

var (
	errMissingEMT = errors.New("missing energy, mass, and/or time")
	errMissingPE  = errors.New("missing password and/or entropy")
)

// populate will solve for entropy, guesses per second, and energy if they aren't given.
// If they are given, it updates them if the computed value is a greater bottleneck than the given value.
func (givens *Givens) populate() error {
	populateDefaults(givens)
	if givens.Password != "" {
		computedEntropy := calculateEntropy(givens.Password)
		if givens.Entropy == 0 || givens.Entropy > computedEntropy {
			givens.Entropy = computedEntropy
		}
	}
	var bremermannGPS float64
	if givens.GuessesPerSecond == 0 && givens.Mass != 0 {
		bremermannGPS = Bremermann * givens.Mass
	}
	calculatePower(givens)
	powerGPS := givens.Power / givens.EnergyPerGuess
	for _, gps := range [2]float64{bremermannGPS, powerGPS} {
		if givens.GuessesPerSecond == 0 || (gps > 0 && gps < givens.GuessesPerSecond) {
			givens.GuessesPerSecond = gps
		}
	}
	calculateEnergy(givens)
	if givens.Energy == 0 && givens.Time == 0 {
		return fmt.Errorf("populating givens: %w", errMissingEMT)
	}
	return nil
}

// BruteForceability computes the liklihood that a password will be
// brute-forced given the contstraints in givens.
// if 0 < BruteForceability <= 1, it represents the probability that the
// password can be brute-forced.
// if BruteForceability > 1, it represents the number of times a password
// can be brute-forced with certainty.
func BruteForceability(givens *Givens, quantum bool) (float64, error) {
	err := givens.populate()
	if err != nil {
		return 0, err
	}
	if givens.Entropy == 0 {
		return 0, fmt.Errorf("BruteForceability: %w", errMissingPE)
	}
	// with Grover's algorithm, quantum computers get an exponential speedup
	var effectiveEntropy float64
	if quantum {
		effectiveEntropy = givens.Entropy / 2
	} else {
		effectiveEntropy = givens.Entropy
	}
	guessesRequired := math.Exp2(effectiveEntropy)
	energyBound := givens.Energy / (guessesRequired * givens.EnergyPerGuess)
	if givens.Time > 0 {
		timeToGuess := guessesRequired / givens.GuessesPerSecond
		timeBound := givens.Time / timeToGuess
		return math.Min(energyBound, timeBound), nil
	}
	return energyBound, nil
}

// MinEntropy calculates the maximum password entropy that the MOAC can certainly brute-force.
// Passwords need an entropy greater than this to have a chance of not being guessed.
func MinEntropy(givens *Givens, quantum bool) (entropy float64, err error) {
	err = givens.populate()
	if err != nil {
		return 0, err
	}
	energyBound := math.Log2(givens.Energy / givens.EnergyPerGuess)
	if givens.Time > 0 {
		timeBound := math.Log2(givens.Time * givens.GuessesPerSecond)
		entropy = math.Min(energyBound, timeBound)
	} else {
		entropy = energyBound
	}

	if quantum {
		return entropy * 2, nil
	}
	return entropy, nil
}
