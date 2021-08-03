package moac

import (
	"testing"
)

var pwgenTests = []struct {
	name           string
	charsetsWanted []string
	entropyWanted  float64
}{
	{
		name:           "everything",
		charsetsWanted: []string{"lowercase", "uppercase", "numbers", "symbols", "latin", "世界🧛"},
		entropyWanted:  256,
	},
	{
		name:           "alnum",
		charsetsWanted: []string{"lowercase", "uppercases", "numbers"},
		entropyWanted:  64,
	},
	{
		name: "tinyPassword",
		charsetsWanted: []string{
			"uppercase", "numbers", "lowercase", "numbers", "numbers", "symbols", "lowercase", "ipaExtensions", "🧛",
		},
		entropyWanted: 1,
	},
	{
		name: "multipleCustomCharsets",
		charsetsWanted: []string{
			"uppercase",
			"numbers",
			"lowercase",
			"𓂸",
			"عظ؆ص",
			"ἀἁἂἃἄἅἆἇἈἉἊἋἌἍἎἏἐἑἒἓἔἕἘἙἚἛἜἝἠἡἢἣἤἥἦἧἨἩἪἫἬἭἮἯἰἱἲἳἴἵἶἷἸἹἺἻἼἽἾἿὀὁὂὃὄὅὈὉὊὋὌὍὐὑὒὓὔὕὖὗὙὛὝὟὠὡὢὣὤὥὦὧὨὩὪὫὬὭὮὯὰάὲέὴήὶίὸόὺύὼώᾀᾁᾂᾃᾄᾅᾆᾇᾈᾉᾊᾋᾌᾍᾎᾏᾐᾑᾒᾓᾔᾕᾖᾗᾘᾙᾚᾛᾜᾝᾞᾟᾠᾡᾢᾣᾤᾥᾦᾧᾨᾩᾪᾫᾬᾭᾮᾯᾰᾱᾲᾳᾴᾶᾷᾸᾹᾺΆᾼ᾽ι᾿῀῁ῂῃῄῆῇῈΈῊΉῌ῍῎῏ῐῑῒΐῖῗῘῙῚΊ῝῞῟ῠῡῢΰῤῥῦῧῨῩῪΎῬ῭΅`ῲῳῴῶῷῸΌῺΏῼ",
		},
		entropyWanted: 256,
	},
}

// second param should include at least one element of the first param.
func latterUsesFormer(former []rune, latter *[]rune) bool {
	for _, char := range former {
		for _, pwChar := range *latter {
			if pwChar == char {
				return true
			}
		}
	}
	return false
}

func pwUsesEachCharset(charsets *[][]rune, password *[]rune) (string, bool) {
	for _, charset := range *charsets {
		if !latterUsesFormer(charset, password) {
			return string(charset), false
		}
	}
	return "", true
}

func pwOnlyUsesAllowedRunes(charsets *[][]rune, password *[]rune) (rune, bool) {
	var allowedChars string
	for _, charset := range *charsets {
		allowedChars += string(charset)
	}
	allowedRunes := []rune(allowedChars)

	charSpace := len(allowedRunes)
	for _, pwChar := range *password {
		for i, char := range allowedRunes {
			if pwChar == char {
				break
			} else if i == charSpace-1 {
				return pwChar, false
			}
		}
	}
	return ' ', true
}

func TestGenPw(t *testing.T) {
	for _, test := range pwgenTests {
		charsetsWanted := &test.charsetsWanted
		entropyWanted := &test.entropyWanted
		t.Run(test.name, func(t *testing.T) {
			charsets := buildCharsets(charsetsWanted)
			// we're dealing with random passwords; try each testcase multiple times
			for i := 0; i < 10; i++ {
				password, err := GenPW(*charsetsWanted, *entropyWanted)
				if err != nil {
					t.Fatalf("GenPW() = %v", err)
				}
				pwRunes := []rune(password)
				if unusedCharset, validPW := pwUsesEachCharset(&charsets, &pwRunes); !validPW {
					t.Errorf("GenPW() = %s; didn't use each charset\nunused charset: %s", password, unusedCharset)
					i = 10
				}
				if invalidRune, validPW := pwOnlyUsesAllowedRunes(&charsets, &pwRunes); !validPW {
					t.Errorf("GenPW() = %s; used invalid character \"%v\"", password, string(invalidRune))
					i = 10
				}
				if entropy := calculateEntropy(password); entropy < *entropyWanted {
					t.Errorf("GenPW() = %s; entropy was %.3g, wanted %.3g", password, entropy, *entropyWanted)
					i = 10
				}
			}
		})
	}
}
