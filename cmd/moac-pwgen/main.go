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
  -G <guesses>	Guesses per second in a brute-force attack
  -l <length>	minimum generated password length; can override (increase) -s
  -L <length>	maximum generated password length; can override (decrease) -s
`
	helpText = "moac-pwgen - generate passwords with the described strength\n" + usage
)

func parseOpts(
	opts *[]getopt.Option, pwr *pwgen.PwRequirements,
) (givens moac.Givens, quantum, exitEarly bool, err error) {
	var (
		minLen64 int64
		maxLen64 int64
	)

	for _, opt := range *opts {
		switch opt.Option {
		case 'h':
			fmt.Fprint(os.Stderr, helpText)

			exitEarly = true

		case 'v':
			fmt.Println(cli.GetVersion())

			exitEarly = true

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
		case 'G':
			givens.GuessesPerSecond, err = strconv.ParseFloat(opt.Value, 64)
		case 'l':
			minLen64, err = strconv.ParseInt(opt.Value, 10, 32)
		case 'L':
			maxLen64, err = strconv.ParseInt(opt.Value, 10, 32)
		}

		if err != nil {
			err = fmt.Errorf("%w: invalid value for -%c: %s\n%s", cli.ErrBadCmdline, opt.Option, opt.Value, usage)

			break
		}

		if exitEarly {
			break
		}
	}

	pwr.MinLen = int(minLen64)
	pwr.MaxLen = int(maxLen64)

	return givens, quantum, exitEarly, err
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
	opts, optind, err := getopt.Getopts(os.Args, "hvqre:s:m:g:P:T:t:G:l:L:")
	if !cli.DisplayErr(err, usage) {
		return 1
	}

	var pwr pwgen.PwRequirements
	givens, quantum, exitEarly, err := parseOpts(&opts, &pwr)

	if !cli.DisplayErr(err, "") {
		return 1
	}

	if exitEarly {
		return 0
	}

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

	charsetNames = []string{"ascii"}

	if len(args) > 0 {
		charsetNames, badCharsets = sanitize.FilterStrings(args)
		warnOnBadCharacters(badCharsets)
	}

	return charsetNames
}
