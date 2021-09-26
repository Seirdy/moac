// Package charsets contains types, functions, and defaults for charsets used in passwords.
package charsets

import (
	"sort"
)

// Charset is the interface implemented by any charset used to build passwords.
// It contains a charset's name as well as string and []rune representations of its glyphs.
type Charset interface {
	String() string
	Runes() []rune
}

// CustomCharset is a charset constructed from an existing rune array.
// It also features methods for mutability, allowing de-duplication and
// sorting of contents.
type CustomCharset []rune

// Runes returns the underlying []rune the CustomCharset is based on.
func (cc CustomCharset) Runes() []rune {
	runesCopy := make([]rune, len(cc))
	for i := range cc {
		runesCopy[i] = cc[i]
	}

	return runesCopy
}

func (cc CustomCharset) String() string {
	return string(cc)
}

func (cc *CustomCharset) sortContents() {
	sort.Slice(*cc, func(i, j int) bool { return (*cc)[i] < (*cc)[j] })
}

// dedupe performs naive rune deduplication on the charset's contents, assuming it has already been sorted.
func (cc *CustomCharset) dedupe() {
	if len(*cc) < 2 {
		return
	}

	n := 1

	for i := 1; i < len(*cc); i++ {
		if (*cc)[i] != (*cc)[i-1] {
			(*cc)[n] = (*cc)[i]
			n++
		}
	}

	*cc = (*cc)[:n]
}

// CharsetCollection holds a list of charsets, and can bulk-convert them to strings, runes, and names.
type CharsetCollection []Charset

// Add adds a charset to a CharsetCollection after de-duplicating its contents and the existing entries.
// All newCharsets are first individually deduplicated/sorted first.
// Whether an existing charset or the new entry gets extra deduplication
// depends on whichever will be bigger afterward; Add will try to
// maximize charset sizes while eliminating any redundancies.
func (cs *CharsetCollection) Add(newCharsets ...Charset) {
	for _, c := range newCharsets {
		cc := CustomCharset(c.Runes())
		cc.sortContents()
		cc.dedupe()

		cs.addSingle(cc)
	}
}

// AddDefault is equivalent to Add, but skips sorting/deduplication.
// This makes adding default charsets a bit faster, since those are already sorted/deduplicated.
func (cs *CharsetCollection) AddDefault(newCharsets ...DefaultCharset) {
	for _, c := range newCharsets {
		cc := CustomCharset(c.Runes())
		cs.addSingle(cc)
	}
}

func (cs *CharsetCollection) addSingle(c CustomCharset) {
	for i := 0; i < len(*cs) && len(c) > 0; i++ {
		minimizeRedundancyInLatter(&(*cs)[i], &c)
	}

	if len(c) > 0 {
		*cs = append(*cs, c)
	}
}

// ParseCharsets creates a CharsetCollection from string identifiers.
// The strings "lowercase", "uppercase", "numbers", and "symbols" all
// refer to their respective constants; "ascii" is an alias for all
// four. The same goes for "latin1", "latinExtendedA", "latinExtendedB",
// and "ipaExtensions"; "latin" is an alias for these four.
// Any other strings have their runes extracted and turned into a new
// CustomCharset.
// See CharsetCollection.Add() for docs on how each entry is added.
func ParseCharsets(charsetNames []string) (cs CharsetCollection) {
	for _, charsetName := range charsetNames {
		switch charsetName {
		case "ascii":
			cs.AddDefault(Lowercase, Uppercase, Numbers, Symbols)
		case "latin":
			cs.AddDefault(Latin1, LatinExtendedA, LatinExtendedB, IPAExtensions)
		default:
			found := false

			for _, defaultCharset := range DefaultCharsets {
				if charsetName == defaultCharset.Name() {
					cs.AddDefault(defaultCharset)

					found = true

					break
				}
			}

			if !found {
				cs.Add(CustomCharset([]rune(charsetName)))
			}
		}
	}

	return cs
}

//nolint:gocritic // c1 is ptr bc it's modified
func minimizeRedundancyInLatter(c1 *Charset, c2 *CustomCharset) {
	c1c := CustomCharset((*c1).Runes())
	moveOverlapToSmaller(&c1c, c2)

	if len(c1c) == 0 {
		*c1, *c2 = *c2, c1c

		return
	}

	*c1 = c1c
}

func moveOverlapToSmaller(c1, c2 *CustomCharset) {
	deleteFromMe := c1
	preserveMe := c2

	if len(*c1) < len(*c2) {
		preserveMe = c1
		deleteFromMe = c2
	}

	for i := 0; i < len(*deleteFromMe); i++ {
		for j := 0; j < len(*preserveMe); j++ {
			if (*deleteFromMe)[i] != (*preserveMe)[j] {
				continue
			}

			*deleteFromMe = append((*deleteFromMe)[:i], (*deleteFromMe)[i+1:]...)
			i--

			break
		}
	}
}
