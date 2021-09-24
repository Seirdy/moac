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
	Name() string
}

// CustomCharset is a charset constructed from an existing rune array.
// This was originally made for moac-pwgen's support for user-provided
// charsets in the form of strings containing desired runes.
type CustomCharset []rune

// Runes returns the original []rune the CustomCharset is based on.
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

// Name is a dummy function for CustomCharsets, since custom charsets don't need name aliases.
// This is because custom charsets are given as a string containing the
// desired runes; we don't have to look them up by name to see what
// runes to fill them with.
func (cc CustomCharset) Name() string {
	return cc.String()
}

// CharsetCollection holds a list of charsets, and can bulk-convert them to strings, runes, and names.
type CharsetCollection []Charset

func (cs *CharsetCollection) addSingle(c CustomCharset) {
	c.sortContents()
	c.dedupe()

	for i := range *cs {
		moveOverlapToSmaller(&(*cs)[i], &c)

		if len((*cs)[i].Runes()) == 0 {
			(*cs)[i] = c

			return
		}
	}

	if len(c) > 0 {
		*cs = append(*cs, c)
	}
}

// Add adds a charset to a CharsetCollection after de-duplicating its contents and the existing entries.
// All newCharsets are first individually deduplicated/sorted first.
// Whether an existing charset or the new entry gets extra deduplication
// depends on whichever will be bigger afterward; Add will try to
// maximize charset sizes while eliminating any redundancies.
func (cs *CharsetCollection) Add(newCharsets ...Charset) {
	for _, c := range newCharsets {
		cs.addSingle(CustomCharset(c.Runes()))
	}
}

func separateOverlap(c1, c2 *CustomCharset) (overlap CustomCharset) {
	for i := 0; i < len(*c1); i++ {
		for j := 0; j < len(*c2); j++ {
			if (*c1)[i] != (*c2)[j] {
				continue
			}

			overlap = append(overlap, (*c1)[i])
			*c1 = append((*c1)[:i], (*c1)[i+1:]...)
			*c2 = append((*c2)[:j], (*c2)[j+1:]...)
			i--

			break
		}
	}

	return overlap
}

func moveOverlapToSmaller(c1 *Charset, c2 *CustomCharset) { //nolint:gocritic // c1 is ptr bc it's modified
	c1c := newCustomCharset(*c1)
	overlap := separateOverlap(&c1c, c2)

	if len(c1c) <= len(*c2) {
		c1c = append(c1c, overlap...)
		c1c.sortContents()

		*c1 = c1c

		return
	}

	*c1 = c1c

	*c2 = append(*c2, overlap...)
	c2.sortContents()
}

func newCustomCharset(c Charset) (fc CustomCharset) {
	return CustomCharset(c.Runes())
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
