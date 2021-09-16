package main

import (
	"fmt"
	"os"
	"strconv"

	"git.sr.ht/~seirdy/moac"
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
) (*moac.Givens, bool, int, int) {
	var (
		givens  moac.Givens
		quantum bool
		minLen  int64
		maxLen  int64
		err     error
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
			givens.Energy, err = strconv.ParseFloat(opt.Value, 64)
		case 's':
			givens.Entropy, err = strconv.ParseFloat(opt.Value, 32)
		case 'm':
			givens.Mass, err = strconv.ParseFloat(opt.Value, 64)
		case 'g':
			givens.EnergyPerGuess, err = strconv.ParseFloat(opt.Value, 64)
		case 'P':
			givens.Power, err = strconv.ParseFloat(opt.Value, 64)
		case 'T':
			givens.Temperature, err = strconv.ParseFloat(opt.Value, 64)
		case 't':
			givens.Time, err = strconv.ParseFloat(opt.Value, 64)
		case 'l':
			minLen, err = strconv.ParseInt(opt.Value, 10, 32)
		case 'L':
			maxLen, err = strconv.ParseInt(opt.Value, 10, 32)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid value for -%c: %s\n%s", opt.Option, opt.Value, helpText)
			os.Exit(1)
		}
	}

	return &givens, quantum, int(minLen), int(maxLen)
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

	var charsets []string

	if len(args) > 0 {
		charsets = args
	} else {
		charsets = []string{"lowercase", "uppercase", "numbers", "symbols", " "}
	}

	pw, err := pwgen.GenPW(charsets, entropyLimit, minLen, maxLen)
	if err != nil {
		fmt.Fprintf(os.Stderr, "moac: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(pw)
}
