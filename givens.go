package main

import (
	"math"

	"github.com/nbutton23/zxcvbn-go"
)

// Givens holds the values used to compute password strength.
// This will grow as the program matures.
// Eventually it'll get its own file and functions to solve for missing vals
// TODO: add power, and use it to compute guesses per second.
// final guesses per second = min(computed, given, Bremermann's limit)
type Givens struct {
	Password         string
	Entropy          float64
	Energy           float64
	Mass             float64 // mass used to build a computer or convert to energy
	Time             float64 // Duration of the attack, in seconds.
	EnergyPerGuess   float64
	guessesPerSecond float64
}

const (
	C         = 299792458      // speed of light in a vacuum, m/s
	G         = 6.67408e-11    // gravitation constant, m^3/kg/s^2
	Hubble    = 2.2e-18        // Hubble's Constant, hertz
	Temp      = 2.7            // cosmic radiation temperature (low estimate), kelvin
	Boltzmann = 1.3806503e-23  // Boltzmann's constant, J/K
	Planck    = 6.62607015e-35 // Planck's Constant, J*s

	// mass of the observable universe
	UMass      = C * C * C / (2 * G * Hubble)
	Bremermann = C * C / Planck              // Bremermann's limit
	Landauer   = Boltzmann * Temp * math.Ln2 // Landauer limit
)

func populateDefaults(givens *Givens) {
	if givens.Mass == 0 {
		// mass of the observable universe
		givens.Mass = UMass
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

// populate will solve for the variables we need to find password strength if they aren't given. If they are given, it updates them if the computed value is a greater bottleneck than the given value.
func (givens *Givens) populate() {
	populateDefaults(givens)
	if givens.Password != "" {
		computedEntropy := calculateEntropy(givens.Password)
		if givens.Entropy == 0 || givens.Entropy > computedEntropy {
			givens.Entropy = computedEntropy
		}
	}
	computedEnergy := givens.Mass * C * C
	if givens.Energy == 0 || givens.Energy > computedEnergy {
		givens.Energy = computedEnergy
	}
	if givens.guessesPerSecond == 0 && givens.Mass != 0 {
		givens.guessesPerSecond = Bremermann * givens.Mass
	}
}

// bruteForceability computes the liklihood that a password will be
// brute-forced given the contstraints in givens.
// if 0 < bruteForceability <= 1, it represents the probability that the
// password can be brute-forced.
// if bruteForceability > 1, it represents the number of times a password
// can be brute-forced with certainty.
func bruteForceability(givens *Givens) float64 {
	givens.populate()
	guessesRequired := math.Pow(2, givens.Entropy)
	energyBound := givens.Energy / (guessesRequired * givens.EnergyPerGuess)
	if givens.Time > 0 {
		timeToGuess := guessesRequired / givens.guessesPerSecond
		timeBound := givens.Time / timeToGuess
		return math.Min(energyBound, timeBound)
	}
	return energyBound
}

// minEntropy calculates the maximum password entropy that the MOAC can certainly brute-force. Passwords need an entropy greater than this to have a chance of not being guessed.
func minEntropy(givens *Givens) float64 {
	givens.populate()
	energyBound := math.Log2(givens.Energy / givens.EnergyPerGuess)
	if givens.Time > 0 {
		timeBound := math.Log2(givens.Time * givens.guessesPerSecond)
		return math.Min(energyBound, timeBound)
	}
	return energyBound
}
