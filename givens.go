package moac

import (
	"errors"
	"fmt"
	"log"
	"math"

	"git.sr.ht/~seirdy/moac/entropy"
	"git.sr.ht/~seirdy/moac/internal/bounds"
)

// Givens holds the "given" values used to compute password strength.
// These values are all physical quantities, measured using standard SI units.
type Givens struct {
	Password         string
	Entropy          float64
	Energy           float64
	Mass             float64 // mass used to build a computer or convert to energy
	Time             float64 // Duration of the attack, in seconds.
	Temperature      float64 // Duration of the attack, in seconds.
	EnergyPerGuess   float64
	Power            float64
	GuessesPerSecond float64
}

const (
	// C is the speed of light in a vacuum, m/s.
	C = 299792458
	// G is the gravitation constant, m^3/kg/s^2.
	G = 6.67408e-11
	// Hubble is Hubble's Constant, hertz.
	Hubble = 2.2e-18
	// UTemp is a low estimate for the temperature of cosmic background radiation, kelvin.
	UTemp = 2.7
	// Boltzmann is Boltzmann's constant, J/K.
	Boltzmann = 1.3806503e-23
	// Planck is Planck's Constant, J*s.
	Planck = 6.62607015e-35

	// UMass is the mass of the observable universe.
	UMass = C * C * C / (2 * G * Hubble)
	// Bremermann is Bremermann's limit.
	Bremermann = C * C / Planck

	// DefaultEntropy is the number of bits of entropy to target if no target entropy is provided.
	DefaultEntropy = 256
)

// landauer outputs the Landauer Limit.
// See https://en.wikipedia.org/wiki/Landauer%27s_principle
func landauer(temp float64) float64 {
	return Boltzmann * temp * math.Ln2
}

// populateDefaults fills in default values for entropy calculation if not provided.
func (givens *Givens) populateDefaults() {
	if givens.Energy+givens.Mass == 0 {
		// mass of the observable universe
		givens.Mass = UMass
	}

	if givens.Entropy == 0 {
		if givens.Mass+givens.EnergyPerGuess == 0 {
			givens.Entropy = DefaultEntropy
		}
	}

	if givens.Temperature == 0 {
		givens.Temperature = UTemp
	}

	if givens.EnergyPerGuess == 0 {
		// maybe put something more elaborate here given different constraints
		givens.EnergyPerGuess = landauer(givens.Temperature)
	}
}

func setBottleneck(computedValues [2]float64, given *float64) {
	for _, computedValue := range computedValues {
		if *given == 0 || (computedValue > 0 && computedValue < *given) {
			*given = computedValue
		}
	}
}

func (givens *Givens) calculatePower() {
	var (
		powerFromComputationSpeed = givens.GuessesPerSecond * givens.EnergyPerGuess
		powerFromEnergy           = givens.Energy / givens.Time
		// loop over an array for this since its length will grow in the future
		computedPowers = [2]float64{powerFromComputationSpeed, powerFromEnergy}
	)

	setBottleneck(computedPowers, &givens.Power)
}

func (givens *Givens) calculateEnergy() {
	var (
		energyFromMass   = givens.Mass * C * C
		energyFromPower  = givens.Power * givens.Time
		computedEnergies = [2]float64{energyFromMass, energyFromPower}
	)

	setBottleneck(computedEnergies, &givens.Energy)
}

// Errors for missing physical values that are required to compute desired values.
var (
	ErrMissingValue = errors.New("not enough given values")
	ErrMissingEMT   = fmt.Errorf("%w: missing energy, mass, and/or time", ErrMissingValue)
	ErrMissingPE    = fmt.Errorf("%w: missing password and/or entropy", ErrMissingValue)
)

// validate ensures that the values in Givens aren't physically impossible.
func (givens *Givens) validate() error {
	if err := bounds.ValidateTemperature(givens.Temperature); err != nil {
		return fmt.Errorf("invalid physical value: %w", err)
	}

	return nil
}

// Populate will solve for entropy, guesses per second, and energy if they aren't given.
// If they are given, it updates them if the computed value is a greater bottleneck than the given value.
func (givens *Givens) Populate() error {
	givens.populateDefaults()

	if err := givens.validate(); err != nil {
		return fmt.Errorf("cannot work with given values: %w", err)
	}

	if givens.Password != "" {
		computedEntropy, err := entropy.Entropy(givens.Password)
		if err != nil {
			log.Panicf("error measuring generated password entropy: %v", err)
		}

		if givens.Entropy == 0 || givens.Entropy > computedEntropy {
			givens.Entropy = computedEntropy
		}
	}

	var bremermannGPS float64

	if givens.GuessesPerSecond == 0 && givens.Mass != 0 {
		bremermannGPS = Bremermann * givens.Mass
	}

	givens.calculatePower()

	powerGPS := givens.Power / givens.EnergyPerGuess
	setBottleneck([2]float64{bremermannGPS, powerGPS}, &givens.GuessesPerSecond)

	givens.calculateEnergy()

	if givens.Energy == 0 && givens.Time == 0 {
		log.Panic("populating givens: failed to populate energy and time")
	}

	return nil
}

// BruteForceability computes the liklihood that a password will be
// brute-forced given the contstraints in givens.
// if 0 < BruteForceability <= 1, it represents the probability that the
// password can be brute-forced.
// if BruteForceability > 1, it represents the number of times a password
// can be brute-forced with certainty.
func (givens *Givens) BruteForceability() (float64, error) {
	if err := givens.Populate(); err != nil {
		return 0, fmt.Errorf("BruteForceability: %w", err)
	}

	if givens.Entropy+givens.Time == 0 {
		return 0, fmt.Errorf("BruteForceability: cannot compute entropy: %w", ErrMissingPE)
	}

	computedBruteForceability := computeBruteForceability(givens)

	// if bruteforceability isn't valid, we have a bug.
	if computedBruteForceability == 0 || math.IsNaN(computedBruteForceability) {
		log.Panicf("failed to compute BruteForceability: got %v", computedBruteForceability)
	}

	return computedBruteForceability, nil
}

// BruteForceabilityQuantum is equivalent to BruteForceability, but accounts for
// quantum computers that use Grover's Algorithm.
func (givens *Givens) BruteForceabilityQuantum() (float64, error) {
	if err := givens.Populate(); err != nil {
		return 0, fmt.Errorf("BruteForceabilityQuantum: %w", err)
	}

	givensQuantum := givens

	// Grover's Algo makes quantum computers as efficient as classical computers at double the entropy.
	givensQuantum.Entropy /= 2

	return givensQuantum.BruteForceability()
}

func computeBruteForceability(givens *Givens) float64 {
	var (
		guessesRequired = math.Exp2(givens.Entropy)
		energyBound     = givens.Energy / (guessesRequired * givens.EnergyPerGuess)
	)

	if givens.Time > 0 {
		timeBound := givens.Time * givens.GuessesPerSecond / guessesRequired

		return math.Min(energyBound, timeBound)
	}

	return energyBound
}

// MinEntropy calculates the maximum password entropy that the MOAC can certainly brute-force.
// Passwords need an entropy greater than this to have a chance of not being guessed.
func (givens *Givens) MinEntropy() (entropyNeeded float64, err error) {
	if err := givens.Populate(); err != nil {
		return 0, fmt.Errorf("MinEntropy: %w", err)
	}

	energyBound := math.Log2(givens.Energy / givens.EnergyPerGuess)

	if givens.Time > 0 {
		timeBound := math.Log2(givens.Time * givens.GuessesPerSecond)
		entropyNeeded = math.Min(energyBound, timeBound)
	} else {
		entropyNeeded = energyBound
	}

	return entropyNeeded, nil
}

// MinEntropyQuantum is equivalent to MinEntropy, but accounts for
// quantum computers that use Grover's Algorithm.
func (givens *Givens) MinEntropyQuantum() (entropyNeeded float64, err error) {
	minEntropyNonQuantum, err := givens.MinEntropy()

	return minEntropyNonQuantum * 2, err
}
