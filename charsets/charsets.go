// Package charsets contains types, functions, and defaults for charsets used in passwords.
package charsets

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
func (CustomCharset) Name() string {
	return ""
}

// CharsetCollection holds a list of charsets, and can bulk-convert them to strings, runes, and names.
type CharsetCollection []Charset

// Strings returns each charset's string representation.
func (cc CharsetCollection) Strings() []string {
	stringArray := make([]string, len(cc))
	for i, c := range cc.Iter() {
		stringArray[i] = c.String()
	}

	return stringArray
}

// Runes returns each charset's rune representation.
func (cc CharsetCollection) Runes() [][]rune {
	runesArray := make([][]rune, len(cc))
	for i, c := range cc {
		runesArray[i] = c.Runes()
	}

	return runesArray
}

// Names returns the human-readable name of each contained charset.
func (cc CharsetCollection) Names() []string {
	namesArray := make([]string, len(cc))
	for i, c := range cc {
		namesArray[i] = c.Name()
	}

	return namesArray
}

// Iter just returns a plain []Charset representation of the CharsetCollection.
func (cc CharsetCollection) Iter() (c []Charset) {
	c = make([]Charset, len(cc))
	for i := range cc {
		c[i] = cc[i]
	}

	return c
}
