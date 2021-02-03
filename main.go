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
  -e <energy>	Maximum energy used by attacker (J).
  -s <entropy>	Password entropy.
  -m <mass>	Mass at attacker's disposal (kg).
  -g <energy>	Energy used per guess (J).
  -t <time>	Time limit for brute-force attack (s).
  -p <password>	Password to analyze.

COMMANDS:
  strength	Calculate the liklihood of a successful guess 
  entropy-limit	Calculate the minimum entropy for a brute-force attack failure.
`

func main() {
	var givens Givens
	opts, optind, err := getopt.Getopts(os.Args, "he:s:m:g:t:p:")
	if err != nil {
		panic(err)
	}
	for _, opt := range opts {
		switch opt.Option {
		case 'h':
			fmt.Println(Usage)
			os.Exit(0)
		case 'e':
			givens.Energy, err = strconv.ParseFloat(opt.Value, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "moac-pwtools: %v", err)
			}
		case 's':
			givens.Entropy, err = strconv.ParseFloat(opt.Value, 32)
			if err != nil {
				fmt.Fprintf(os.Stderr, "moac-pwtools: %v", err)
			}
		case 'm':
			givens.Mass, err = strconv.ParseFloat(opt.Value, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "moac-pwtools: %v", err)
			}
		case 'g':
			givens.EnergyPerGuess, err = strconv.ParseFloat(opt.Value, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "moac-pwtools: %v", err)
			}
		case 't':
			givens.Time, err = strconv.ParseFloat(opt.Value, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "moac-pwtools: %v", err)
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
			fmt.Printf("%.3g\n", bruteForceability(&givens))
		case "entropy-limit":
			fmt.Printf("%.3g\n", minEntropy(&givens))
		default:
			log.Println(Usage)
			os.Exit(1)
		}
	} else {
		fmt.Printf("%.3g\n", bruteForceability(&givens))
	}
}
