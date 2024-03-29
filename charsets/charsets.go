// Package charsets contains types, functions, and defaults for charsets used in passwords.
package charsets

import (
	"sort"
	"strings"
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

	removals := 1

	for i := 1; i < len(*cc); i++ {
		if (*cc)[i] != (*cc)[i-1] {
			(*cc)[removals] = (*cc)[i]
			removals++
		}
	}

	*cc = (*cc)[:removals]
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
	for i := 0; i < len(*cs); i++ {
		minimizeRedundancyInLatter(&(*cs)[i], &c)
	}

	if len(c) > 0 {
		*cs = append(*cs, c)
	}
}

// Combined returns a CustomCharset that combines all the charsets in a CharsetCollection.
func (cs *CharsetCollection) Combined() CustomCharset {
	var ccBuilder strings.Builder

	for _, c := range *cs {
		ccBuilder.WriteString(c.String())
	}

	return []rune(ccBuilder.String())
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
			parseCharset(&cs, charsetName)
		}
	}

	return cs
}

func parseCharset(cs *CharsetCollection, charsetName string) {
	found := false

	for _, defaultCharset := range DefaultCharsets {
		if charsetName != defaultCharset.Name() {
			continue
		}

		cs.AddDefault(defaultCharset)

		found = true

		break
	}

	if !found {
		cs.Add(CustomCharset([]rune(charsetName)))
	}
}

//nolint:gocritic // c1 is ptr bc it's modified
func minimizeRedundancyInLatter(former *Charset, latter *CustomCharset) {
	c1c := CustomCharset((*former).Runes())
	moveOverlapToSmaller(&c1c, latter)

	if len(c1c) == 0 {
		*former, *latter = *latter, c1c

		return
	}

	*former = c1c
}

func moveOverlapToSmaller(c1, c2 *CustomCharset) {
	deleteFromMe := c1
	preserveMe := c2

	if len(*c1) < len(*c2) {
		preserveMe = c1
		deleteFromMe = c2
	}

	for delI := 0; delI < len(*deleteFromMe); delI++ {
		for presI := 0; presI < len(*preserveMe); presI++ {
			if (*deleteFromMe)[delI] != (*preserveMe)[presI] {
				continue
			}

			*deleteFromMe = append((*deleteFromMe)[:delI], (*deleteFromMe)[delI+1:]...)
			delI--

			break
		}
	}
}
