package main

import (
	"fmt"
	"os"
	"strconv"

	"git.sr.ht/~sircmpwn/getopt"
)

const Usage = `moac-pwtools - analyze password strength with physical limits

USAGE:
  moac-pwtools [OPTIONS] [COMMAND] [ARGS]

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

func main() {
	var (
		givens  Givens
		quantum bool
	)
	opts, optind, err := getopt.Getopts(os.Args, "hqe:s:m:g:P:t:p:")
	if err != nil {
		fmt.Fprintf(os.Stderr, "moac-pwtools: %v\n%s", err, Usage)
		os.Exit(1)
	}
	for _, opt := range opts {
		switch opt.Option {
		case 'h':
			fmt.Println(Usage)
			os.Exit(0)
		case 'q':
			quantum = true
		case 'e':
			givens.Energy, err = strconv.ParseFloat(opt.Value, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "moac-pwtools: %v\n", err)
				os.Exit(1)
			}
		case 's':
			givens.Entropy, err = strconv.ParseFloat(opt.Value, 32)
			if err != nil {
				fmt.Fprintf(os.Stderr, "moac-pwtools: %v\n%s", err, Usage)
				os.Exit(1)
			}
		case 'm':
			givens.Mass, err = strconv.ParseFloat(opt.Value, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "moac-pwtools: %v\n%s", err, Usage)
				os.Exit(1)
			}
		case 'g':
			givens.EnergyPerGuess, err = strconv.ParseFloat(opt.Value, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "moac-pwtools: %v\n%s", err, Usage)
				os.Exit(1)
			}
		case 'P':
			givens.Power, err = strconv.ParseFloat(opt.Value, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "moac-pwtools: %v\n%s", err, Usage)
				os.Exit(1)
			}
		case 't':
			givens.Time, err = strconv.ParseFloat(opt.Value, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "moac-pwtools: %v\n%s", err, Usage)
				os.Exit(1)
			}
		case 'p':
			givens.Password = opt.Value
		}
	}
	args := os.Args[optind:]
	if len(args) >= 1 {
		cmd := args[0]
		switch cmd {
		case "strength":
			likelihood, err := BruteForceability(&givens, quantum)
			if err != nil {
				fmt.Fprintf(os.Stderr, "moac-pwtools: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("%.3g\n", likelihood)
		case "entropy-limit":
			entropyLimit, err := MinEntropy(&givens, quantum)
			if err != nil {
				fmt.Fprintf(os.Stderr, "moac-pwtools: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("%.3g\n", entropyLimit)
		case "pwgen":
			entropyLimit, err := MinEntropy(&givens, quantum)
			if err != nil {
				fmt.Fprintf(os.Stderr, "moac-pwtools: %v\n", err)
				os.Exit(1)
			}
			var charsets []string
			if len(args) > 1 {
				charsets = args[1:]
			} else {
				charsets = []string{"lowercase", "uppercase", "numbers", "symbols", "extendedASCII", " "}
			}
			pw, err := GenPW(charsets, entropyLimit)
			if err != nil {
				fmt.Fprintf(os.Stderr, "moac-pwtools: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(pw)
		default:
			fmt.Fprintln(os.Stderr, Usage)
			os.Exit(1)
		}
	} else {
		likelihood, err := BruteForceability(&givens, quantum)
		if err != nil {
			fmt.Fprintf(os.Stderr, "moac-pwtools: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%.3g\n", likelihood)
	}
}
