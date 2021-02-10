package main

import (
	"fmt"
	"os"
	"strconv"

	"git.sr.ht/~seirdy/moac"
	"git.sr.ht/~sircmpwn/getopt"
)

const Usage = `moac - analyze password strength with physical limits

USAGE:
  moac [OPTIONS] [COMMAND] [ARGS]

OPTIONS:
  -h	Display this help message.
  -q	Account for quantum computers using Grover's algorithm
  -e <energy>	Maximum energy used by attacker (J).
  -s <entropy>	Password entropy.
  -m <mass>	Mass at attacker's disposal (kg).
  -g <energy>	Energy used per guess (J).
  -P <power>	Power available to the computer (W)
  -t <time>	Time limit for brute-force attack (s).
  -p <password>	Password to analyze.

COMMANDS:
  strength	Calculate the liklihood of a successful guess 
  entropy-limit	Calculate the minimum entropy for a brute-force attack failure.
  pwgen	generate a password resistant to the described brute-force attack,
        using charsets specified by [ARGS] (defaults to all provided charsets)
`

func parseOpts(opts *[]getopt.Option) (moac.Givens, bool, error) {
	var (
		givens  moac.Givens
		quantum bool
		err     error
	)
	for _, opt := range *opts {
		switch opt.Option {
		case 'h':
			fmt.Println(Usage)
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
		case 't':
			givens.Time, err = strconv.ParseFloat(opt.Value, 64)
		case 'p':
			givens.Password = opt.Value
		}
		if err != nil {
			return givens, quantum, fmt.Errorf("bad value for -%c: %w", opt.Option, err)
		}
	}
	return givens, quantum, nil
}

func main() {
	opts, optind, err := getopt.Getopts(os.Args, "hqe:s:m:g:P:t:p:")
	if err != nil {
		fmt.Fprintf(os.Stderr, "moac: %v\n%s", err, Usage)
		os.Exit(1)
	}
	givens, quantum, err := parseOpts(&opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "moac: %v\n%s", err, Usage)
		os.Exit(1)
	}
	args := os.Args[optind:]
	if len(args) == 0 {
		likelihood, err := moac.BruteForceability(&givens, quantum)
		if err != nil {
			fmt.Fprintf(os.Stderr, "moac: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%.3g\n", likelihood)
		os.Exit(0)
	}
	cmd := args[0]
	switch cmd {
	case "strength":
		likelihood, err := moac.BruteForceability(&givens, quantum)
		if err != nil {
			fmt.Fprintf(os.Stderr, "moac: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%.3g\n", likelihood)
	case "entropy-limit":
		entropyLimit, err := moac.MinEntropy(&givens, quantum)
		if err != nil {
			fmt.Fprintf(os.Stderr, "moac: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%.3g\n", entropyLimit)
	case "pwgen":
		entropyLimit, err := moac.MinEntropy(&givens, quantum)
		if err != nil {
			fmt.Fprintf(os.Stderr, "moac: %v\n", err)
			os.Exit(1)
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
		fmt.Fprintln(os.Stderr, Usage)
		os.Exit(1)
	}
}
