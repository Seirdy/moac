package pwgen

// indexing contains a bunch of utilities for working with slice indexes:
// finding indexes, inserting at an index, etc.

import (
	"crypto/rand"
	"log"
	"math/big"
	"strings"
)

func randInt(max int) int {
	newInt, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		log.Panicf("can't generate passwords: crypto/rand unavailable: %v", err)
	}

	return int(newInt.Int64())
}

func addRuneToEnd(password *strings.Builder, runes []rune) {
	newChar := runes[randInt(len(runes))]
	password.WriteRune(newChar)
}

func addRuneAtRandLoc(pwRunes *[]rune, runesToPickFrom []rune) {
	newChar := runesToPickFrom[randInt(len(runesToPickFrom))]
	index := randInt(len(*pwRunes))
	*pwRunes = append((*pwRunes)[:index+1], (*pwRunes)[index:]...)
	(*pwRunes)[index] = newChar
}

func indexOf(src []int, e int) int {
	for i, a := range src {
		if a == e {
			return i
		}
	}

	return -1
}
