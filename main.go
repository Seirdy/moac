package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"git.sr.ht/~sircmpwn/getopt"
)

const Usage = `moac-pwtools - analyze password strength with physical limits

USAGE:
  moac-pwtools [OPTIONS] [COMMAND]

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
`

func main() {
	var (
		givens  Givens
		quantum bool
	)
	opts, optind, err := getopt.Getopts(os.Args, "hqe:s:m:g:P:t:p:")
	if err != nil {
		panic(err)
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
				fmt.Fprintf(os.Stderr, "moac-pwtools: %v\n", err)
				os.Exit(1)
			}
		case 'm':
			givens.Mass, err = strconv.ParseFloat(opt.Value, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "moac-pwtools: %v\n", err)
				os.Exit(1)
			}
		case 'g':
			givens.EnergyPerGuess, err = strconv.ParseFloat(opt.Value, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "moac-pwtools: %v\n", err)
				os.Exit(1)
			}
		case 'P':
			givens.Power, err = strconv.ParseFloat(opt.Value, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "moac-pwtools: %v\n", err)
				os.Exit(1)
			}
		case 't':
			givens.Time, err = strconv.ParseFloat(opt.Value, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "moac-pwtools: %v\n", err)
				os.Exit(1)
			}
		case 'p':
			givens.Password = opt.Value
		}
	}
	args := os.Args[optind:]
	if len(args) > 1 {
		log.Println(Usage)
		os.Exit(1)
	} else if len(args) == 1 {
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
		default:
			log.Println(Usage)
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
