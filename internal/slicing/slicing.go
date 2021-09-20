// Package slicing implements utility functions for slices of runes
package slicing

import (
	"reflect"
)

// MapToSlice extracts the values of original to a 2D rune array.
func MapToSlice(original map[string][]rune) (res [][]rune) {
	for _, value := range original {
		res = append(res, value)
	}

	return res
}

// SliceContainsRuneSlice determines of a 2D rune array contains a 1D rune array.
func SliceContainsRuneSlice(runes [][]rune, item []rune) bool {
	for _, possibleMatch := range runes {
		if reflect.DeepEqual(item, possibleMatch) {
			return true
		}
	}

	return false
}
