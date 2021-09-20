package pwgen

// indexing contains a bunch of utilities for working with slice indexes:
// finding indexes, inserting at an index, etc.

import (
	"crypto/rand"
	"log"
	"math/big"
	"sort"
	"strings"
)

func randInt(max int) int {
	newInt, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		log.Panicf("specialIndexes: %v", err)
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

func removeIndex(s []rune, index int) []rune {
	ret := append(make([]rune, 0), s[:index]...)

	return append(ret, s[index+1:]...)
}

func sortRunes(runes *[]rune) {
	sort.Slice(*runes, func(i, j int) bool { return (*runes)[i] < (*runes)[j] })
}

func dedupeRunes(runes []rune) []rune {
	sortRunes(&runes)

	keys := make(map[rune]bool)
	list := []rune{}

	for _, entry := range runes {
		if _, value := keys[entry]; !value {
			keys[entry] = true

			list = append(list, entry)
		}
	}

	return list
}

// removeLatterFromFormer removes one of each of latter's elements from former and returns the result.
func removeLatterFromFormer(former, latter []rune) (res, overlap []rune) {
	overlap = make([]rune, 0)
	res = make([]rune, len(former))
	copy(res, former)

	for _, latterItem := range latter {
		for i := 0; i < len(res); i++ {
			if res[i] == latterItem {
				// overlap is initialized with length zero, but makezero can't tell.
				overlap = append(overlap, res[i]) //nolint:makezero // false-positive
				res = removeIndex(res, i)
				i--
			}
		}
	}

	return res, overlap
}
