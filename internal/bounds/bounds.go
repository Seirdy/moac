// Package bounds holds/checks max/min acceptable values of user-provided inputs
package bounds

import (
	"errors"
	"fmt"
)

const (
	// PlanckTemp is the Planck Temperature, the maximum temperature that today's physics can handle
	// Temperature values above this will result in an error.
	PlanckTemp = 1.417e32
)

var (
	// ErrImpossiblyHigh indicates that a value is higher than our current
	// understanding of physics allows or accounts for.
	ErrImpossiblyHigh = errors.New("value is physically impossibly large")
	// ErrImpossibleNegative indicates that a value that must be above zero was too low.
	ErrImpossibleNegative = errors.New("value must be above 0")
)

// ValidateTemperature ensures that a given temperature is physically possible.
// Temperatures must be above zero and cannot surpass the Planck Temperature.
func ValidateTemperature(temp float64) error {
	if temp <= 0 {
		return fmt.Errorf("temperature too low: %w", ErrImpossibleNegative)
	}

	if temp > PlanckTemp {
		return fmt.Errorf("temperature above Planck Temperature: %w", ErrImpossiblyHigh)
	}

	return nil
}

// NonNegative validates that all the given values are at or above 0.
func NonNegative(vs ...float64) error {
	for _, v := range vs {
		if v < 0 {
			return fmt.Errorf("physical value too low: %w", ErrImpossibleNegative)
		}
	}

	return nil
}
