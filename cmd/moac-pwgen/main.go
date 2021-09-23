package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"git.sr.ht/~seirdy/moac"
	"git.sr.ht/~seirdy/moac/internal/sanitize"
	"git.sr.ht/~seirdy/moac/internal/version"
	"git.sr.ht/~seirdy/moac/pwgen"
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
	helpText = "moac-pwgen - generate passwords with the described strength" + usage
)

func parseOpts( //nolint:cyclop // complexity solely determined by cli flag count
	opts *[]getopt.Option,
) (givens *moac.Givens, quantum bool, minLen, maxLen int) {
	var (
		givensValue moac.Givens
		minLen64    int64
		maxLen64    int64
		err         error
	)

	for _, opt := range *opts {
		switch opt.Option {
		case 'h':
			fmt.Println(helpText)
			os.Exit(0)
		case 'v':
			fmt.Println(version.GetVersion())
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

	return &givensValue, quantum, int(minLen64), int(maxLen64)
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
	opts, optind, err := getopt.Getopts(os.Args, "hvqre:s:m:g:P:T:t:l:L:")
	if err != nil {
		fmt.Fprintf(os.Stderr, "moac: %v\n%s", err, usage)
		os.Exit(1)
	}

	givens, quantum, minLen, maxLen := parseOpts(&opts)
	args := os.Args[optind:]
	// If the only user-supplied given is entropy, then just use that
	// entropy level and skip calculating the strength of a brute-force attack.
	entropyLimit := givens.Entropy

	if givens.Energy+givens.Mass+givens.Power+givens.Time != 0 {
		entropyLimit = moac.MinEntropy(givens, quantum)
	}

	var charsets, badCharsets []string

	if len(args) > 0 {
		charsets, badCharsets = sanitize.FilterStrings(args)
		warnOnBadCharacters(badCharsets)
	} else {
		charsets = []string{"ascii"}
	}

	pw, err := pwgen.GenPW(charsets, entropyLimit, minLen, maxLen)
	if err != nil {
		fmt.Fprintf(os.Stderr, "moac: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(pw)
}
