package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"git.sr.ht/~seirdy/moac/v2"
	"git.sr.ht/~seirdy/moac/v2/charsets"
	"git.sr.ht/~seirdy/moac/v2/internal/cli"
	"git.sr.ht/~seirdy/moac/v2/internal/sanitize"
	"git.sr.ht/~seirdy/moac/v2/pwgen"
	"git.sr.ht/~sircmpwn/getopt"
)

const (
	usage = `
USAGE:
  moac-pwgen [OPTIONS] [CHARSETS]

OPTIONS:
  -h	Display this help message
  -q	Account for quantum computers using Grover's algorithm
  -e <energy>	Maximum energy used by attacker (J)
  -s <entropy>	Password entropy
  -m <mass>	Mass at attacker's disposal (kg)
  -g <energy>	Energy used per guess (J)
  -P <power>	Power available to the computer (W)
  -T <temperature>	Temperature of the system (K)
  -t <time>	Time limit for brute-force attack (s)
  -l <length>	minimum generated password length; can override (increase) -s
  -L <length>	maximum generated password length; can override (decrease) -s
`
	helpText = "moac-pwgen - generate passwords with the described strength\n" + usage
)

func parseOpts( //nolint:cyclop // complexity solely determined by cli flag count
	opts *[]getopt.Option, pwr *pwgen.PwRequirements,
) (givens *moac.Givens, quantum bool) {
	var (
		givensValue moac.Givens
		minLen64    int64
		maxLen64    int64
		err         error
	)

	for _, opt := range *opts {
		switch opt.Option {
		case 'h':
			fmt.Fprint(os.Stderr, helpText)
			os.Exit(0)
		case 'v':
			fmt.Println(cli.GetVersion())
			os.Exit(0)
		case 'q':
			quantum = true
		case 'e':
			givensValue.Energy, err = strconv.ParseFloat(opt.Value, 64)
		case 's':
			givensValue.Entropy, err = strconv.ParseFloat(opt.Value, 32)
		case 'm':
			givensValue.Mass, err = strconv.ParseFloat(opt.Value, 64)
		case 'g':
			givensValue.EnergyPerGuess, err = strconv.ParseFloat(opt.Value, 64)
		case 'P':
			givensValue.Power, err = strconv.ParseFloat(opt.Value, 64)
		case 'T':
			givensValue.Temperature, err = strconv.ParseFloat(opt.Value, 64)
		case 't':
			givensValue.Time, err = strconv.ParseFloat(opt.Value, 64)
		case 'l':
			minLen64, err = strconv.ParseInt(opt.Value, 10, 32)
		case 'L':
			maxLen64, err = strconv.ParseInt(opt.Value, 10, 32)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid value for -%c: %s\n%s", opt.Option, opt.Value, helpText)
			os.Exit(1)
		}
	}

	pwr.MinLen = int(minLen64)
	pwr.MaxLen = int(maxLen64)

	return &givensValue, quantum
}

func warnOnBadCharacters(badCharsets []string) {
	if len(badCharsets) == 0 {
		return
	}

	var warningSubstring strings.Builder

	for i := 0; i < len(badCharsets)-1; i++ {
		warningSubstring.WriteString(badCharsets[i])
		warningSubstring.WriteString(", ")
	}

	warningSubstring.WriteString(badCharsets[len(badCharsets)-1])

	fmt.Fprintf(os.Stderr, "warning: charsets %v contained invalid codepoints, removing them\n", warningSubstring.String())
}

func main() {
	os.Exit(main1())
}

func main1() int {
	opts, optind, err := getopt.Getopts(os.Args, "hvqre:s:m:g:P:T:t:l:L:")
	if !cli.DisplayErr(err, usage) {
		return 1
	}

	var pwr pwgen.PwRequirements
	givens, quantum := parseOpts(&opts, &pwr)
	args := os.Args[optind:]
	// If the only user-supplied given is entropy, then just use that
	// entropy level and skip calculating the strength of a brute-force attack.
	pwr.TargetEntropy = givens.Entropy

	if givens.Energy+givens.Mass+givens.Power+givens.Time != 0 {
		if quantum {
			pwr.TargetEntropy, err = givens.MinEntropyQuantum()
		} else {
			pwr.TargetEntropy, err = givens.MinEntropy()
		}
	}

	if !cli.DisplayErr(err, "") {
		return 1
	}

	pwr.CharsetsWanted = charsets.ParseCharsets(setCharsetNames(args))
	pw, err := pwgen.GenPW(pwr)

	if !cli.DisplayErr(err, "") {
		return 1
	}

	fmt.Print(pw)

	return 0
}

func setCharsetNames(args []string) (charsetNames []string) {
	var badCharsets []string

	if len(args) > 0 {
		charsetNames, badCharsets = sanitize.FilterStrings(args)
		warnOnBadCharacters(badCharsets)
	} else {
		charsetNames = []string{"ascii"}
	}

	return charsetNames
}
