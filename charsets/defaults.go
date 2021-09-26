package charsets

// DefaultCharset is a pre-built charset, with a name and string/rune representations.
type DefaultCharset int

const (
	// Lowercase contains runes in the range [a-z].
	Lowercase DefaultCharset = iota
	// Uppercase contains runes in the range [A-Z].
	Uppercase
	// Numbers contains runes in the range [0-9].
	Numbers
	// Symbols contains all the ASCII symbols.
	Symbols
	// Latin1 contains all the glyphs in the Latin-1 Unicode block.
	Latin1
	// LatinExtendedA contains all the glyphs in the Latin Extended-A Unicode block.
	LatinExtendedA
	// LatinExtendedB contains all the glyphs in the Latin Extended-B Unicode block.
	LatinExtendedB
	// IPAExtensions contains all the glyphs in the IPA Extensions Unicode block.
	IPAExtensions
)

func (dc DefaultCharset) String() string {
	return [...]string{
		"abcdefghijklmnopqrstuvwxyz",
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"0123456789",
		"!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~",
		"¡¢£¤¥¦§¨©ª«¬®¯°±²³´μ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþÿ",
		"ĀāĂăĄąĆćĈĉĊċČčĎďĐđĒēĔĕĖėĘęĚěĜĝĞğĠġĢģĤĥĦħĨĩĪīĬĭĮįİıĲĳĴĵĶķĸĹĺĻļĽľĿŀŁłŃńŅņŇňŉŊŋŌōŎŏŐőŒœŔŕŖŗŘřŚśŜŝŞşŠšŢţŤťŦŧŨũŪūŬŭŮůŰűŲųŴŵŶŷŸŹźŻżŽžſ",
		"ƀƁƂƃƄƅƆƇƈƉƊƋƌƍƎƏƐƑƒƓƔƕƖƗƘƙƚƛƜƝƞƟƠơƢƣƤƥƦƧƨƩƪƫƬƭƮƯưƱƲƳƴƵƶƷƸƹƺƻƼƽƾƿǀǁǂǃǄǅǆǇǈǉǊǋǌǍǎǏǐǑǒǓǔǕǖǗǘǙǚǛǜǝǞǟǠǡǢǣǤǥǦǧǨǩǪǫǬǭǮǯǰǱǲǳǴǵǶǷǸǹǺǻǼǽǾǿȀȁȂȃȄȅȆȇȈȉȊȋȌȍȎȏȐȑȒȓȔȕȖȗȘșȚțȜȝȞȟȠȡȢȣȤȥȦȧȨȩȪȫȬȭȮȯȰȱȲȳȴȵȶȷȸȹȺȻȼȽȾȿɀɁɂɃɄɅɆɇɈɉɊɋɌɍɎɏ",
		"ɐɑɒɓɔɕɖɗɘəɚɛɜɝɞɟɠɡɢɣɤɥɦɧɨɩɪɫɬɭɮɯɰɱɲɳɴɵɶɷɸɹɺɻɼɽɾɿʀʁʂʃʄʅʆʇʈʉʊʋʌʍʎʏʐʑʒʓʔʕʖʗʘʙʚʛʜʝʞʟʠʡʢʣʤʥʦʧʨʩʪʫʬʭʮʯ",
	}[dc]
}

// Name outputs the human-readable name for a charset, useful for matching user input.
func (dc DefaultCharset) Name() string {
	return [...]string{
		"lowercase",
		"uppercase",
		"numbers",
		"symbols",
		"latin1",
		"latinExtendedA",
		"latinExtendedB",
		"ipaExtensions",
	}[dc]
}

// Runes outputs all the runes in a charset.
func (dc DefaultCharset) Runes() []rune {
	return []rune(dc.String())
}

// DefaultCharsets contains pre-built named charsets.
// This makes identifying a charset by name (e.g. as a CLI argument) much easier.
var DefaultCharsets = []DefaultCharset{
	Lowercase, Uppercase, Numbers, Symbols,
	Latin1, LatinExtendedA, LatinExtendedB, IPAExtensions,
}
