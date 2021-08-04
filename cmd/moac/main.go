package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"

	"git.sr.ht/~seirdy/moac"
	"git.sr.ht/~sircmpwn/getopt"
	"golang.org/x/term"
)

const (
	Usage = `
USAGE:
  moac [OPTIONS] [COMMAND] [ARGS]

OPTIONS:
  -h	Display this help message.
  -q	Account for quantum computers using Grover's algorithm
  -r	Interactively enter a password in the terminal; overrides -p
  -e <energy>	Maximum energy used by attacker (J).
  -s <entropy>	Password entropy.
  -m <mass>	Mass at attacker's disposal (kg).
  -g <energy>	Energy used per guess (J).
  -P <power>	Power available to the computer (W)
  -t <time>	Time limit for brute-force attack (s).
  -p <password>	Password to analyze (do not use a real password).

COMMANDS:
  strength	Calculate the liklihood of a successful guess 
  entropy-limit	Calculate the minimum entropy for a brute-force attack failure.
  pwgen	generate a password resistant to the described brute-force attack,
       	using charsets specified by [ARGS] (defaults to all provided charsets)
`
	helpText = "moac - analyze password strength with physical limits" + Usage
)

func parseOpts(opts *[]getopt.Option) (*moac.Givens, bool, bool, error) {
	var (
		givens       moac.Givens
		quantum      bool
		readPassword bool
		err          error
	)
	for _, opt := range *opts {
		switch opt.Option {
		case 'h':
			fmt.Println(helpText)
			os.Exit(0)
		case 'q':
			quantum = true
		case 'r':
			readPassword = true
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
		case 't':
			givens.Time, err = strconv.ParseFloat(opt.Value, 64)
		case 'p':
			givens.Password = opt.Value
		}
		if err != nil {
			return &givens, quantum, readPassword, fmt.Errorf("invalid value for -%c: %s", opt.Option, opt.Value)
		}
	}
	return &givens, quantum, readPassword, nil
}

func getBruteForceability(givens *moac.Givens, quantum bool) float64 {
	likelihood, err := moac.BruteForceability(givens, quantum)
	if err != nil {
		fmt.Fprintf(os.Stderr, "moac: %v\n", err)
		os.Exit(1)
	}
	return likelihood
}

func getMinEntropy(givens *moac.Givens, quantum bool) float64 {
	entropyLimit, err := moac.MinEntropy(givens, quantum)
	if err != nil {
		fmt.Fprintf(os.Stderr, "moac: %v\n", err)
		os.Exit(1)
	}
	return entropyLimit
}

func main() {
	opts, optind, err := getopt.Getopts(os.Args, "hqre:s:m:g:P:t:p:")
	if err != nil {
		fmt.Fprintf(os.Stderr, "moac: %v\n%s", err, Usage)
		os.Exit(1)
	}
	givens, quantum, readPassword, err := parseOpts(&opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "moac: %v\n%s", err, Usage)
		os.Exit(1)
	}
	if readPassword {
		fmt.Print("Enter password: ")
		bytepw, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			os.Exit(1)
		}
		givens.Password = string(bytepw)
	}
	args := os.Args[optind:]
	if len(args) == 0 {
		fmt.Printf("%.3g\n", getBruteForceability(givens, quantum))
		os.Exit(0)
	}
	cmd := args[0]
	switch cmd {
	case "strength":
		fmt.Printf("%.3g\n", getBruteForceability(givens, quantum))
	case "entropy-limit":
		fmt.Printf("%.3g\n", getMinEntropy(givens, quantum))
	case "pwgen":
		// If the only user-supplied given is entropy, then just use that
		// entropy level and skip calculating the strength of a brute-force attack.
		entropyLimit := givens.Entropy
		if givens.Energy+givens.Mass+givens.Power+givens.Time != 0 {
			entropyLimit = getMinEntropy(givens, quantum)
		}
		var charsets []string
		if len(args) > 1 {
			charsets = args[1:]
		} else {
			charsets = []string{"lowercase", "uppercase", "numbers", "symbols", "extendedASCII", " "}
		}
		pw, err := moac.GenPW(charsets, entropyLimit)
		if err != nil {
			fmt.Fprintf(os.Stderr, "moac: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(pw)
	default:
		fmt.Fprintf(os.Stderr, "moac: unknown command %v\n%s", cmd, Usage)
		os.Exit(1)
	}
}
