package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"

	"git.sr.ht/~seirdy/moac/v2"
	"git.sr.ht/~seirdy/moac/v2/entropy"
	"git.sr.ht/~seirdy/moac/v2/internal/cli"
	"git.sr.ht/~seirdy/moac/v2/internal/version"
	"git.sr.ht/~sircmpwn/getopt"
	"golang.org/x/term"
)

const (
	usage = `
USAGE:
  moac [OPTIONS] [COMMAND]

OPTIONS:
  -h	Display this help message
  -q	Account for quantum computers using Grover's algorithm
  -r	Interactively enter a password in the terminal; overrides -p
  -e <energy>	Maximum energy used by attacker (J)
  -s <entropy>	Password entropy
  -m <mass>	Mass at attacker's disposal (kg)
  -g <energy>	Energy used per guess (J)
  -P <power>	Power available to the computer (W)
  -T <temperature>	Temperature of the system (K)
  -t <time>	Time limit for brute-force attack (s)
  -p <password>	Password to analyze; use "-" for stdin

COMMANDS:
  strength	Calculate the liklihood of a successful guess 
  entropy	Calculate the entropy of the given password
  entropy-limit	Calculate the minimum entropy for a brute-force attack failure.
`
	helpText = "moac - analyze password strength with physical limits" + usage
)

func parseOpts( //nolint:cyclop // complexity solely determined by cli flag count
	opts *[]getopt.Option,
) (*moac.Givens, bool, bool) {
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
		case 'v':
			fmt.Println(version.GetVersion())
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
		case 'T':
			givens.Temperature, err = strconv.ParseFloat(opt.Value, 64)
		case 't':
			givens.Time, err = strconv.ParseFloat(opt.Value, 64)
		case 'p':
			givens.Password = opt.Value
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid value for -%c: %s\n%s", opt.Option, opt.Value, helpText)
			os.Exit(1)
		}
	}

	return &givens, quantum, readPassword
}

func getBruteForceability(givens *moac.Givens, quantum bool) (likelihood float64) {
	var err error

	if quantum {
		likelihood, err = givens.BruteForceabilityQuantum()
	} else {
		likelihood, err = givens.BruteForceability()
	}

	if err != nil {
		cli.ExitOnErr(err, "")
	}

	return likelihood
}

func getEntropy(givens *moac.Givens) float64 {
	if givens.Password == "" {
		fmt.Fprintf(os.Stderr, "moac: cannot compute entropy: missing password\n")
		os.Exit(1)
	}

	return entropy.Entropy(givens.Password)
}

func readPwInteractive(password *string) {
	fmt.Print("Enter password: ")

	bytepw, err := term.ReadPassword(int(syscall.Stdin)) //nolint:unconvert // needed for some platforms

	fmt.Println()

	cli.ExitOnErr(err, "failed to read password")

	*password = string(bytepw)
}

func readPwStdin(password *string) {
	stdinBytes, err := ioutil.ReadAll(os.Stdin)
	cli.ExitOnErr(err, "")

	*password = string(stdinBytes)
}

func getMinEntropy(givens *moac.Givens, quantum bool) float64 {
	var (
		minEntropy float64
		err        error
	)

	if quantum {
		minEntropy, err = givens.MinEntropyQuantum()
	} else {
		minEntropy, err = givens.MinEntropy()
	}

	cli.ExitOnErr(err, "")

	return minEntropy
}

func main() {
	opts, optind, err := getopt.Getopts(os.Args, "hvqre:s:m:g:P:T:t:p:")
	cli.ExitOnErr(err, usage)

	givens, quantum, readPassword := parseOpts(&opts)
	if readPassword {
		readPwInteractive(&givens.Password)
	} else if givens.Password == "-" {
		readPwStdin(&givens.Password)
	}

	givens.Password = strings.TrimSuffix(givens.Password, "\n")
	args := os.Args[optind:]
	cmd := "strength"

	if len(args) > 0 {
		cmd = args[0]
	}

	switch cmd {
	case "strength":
		fmt.Printf(cli.FloatFmt, getBruteForceability(givens, quantum))
	case "entropy":
		fmt.Printf(cli.FloatFmt, getEntropy(givens))
	case "entropy-limit":
		fmt.Printf(cli.FloatFmt, getMinEntropy(givens, quantum))
	default:
		fmt.Fprintf(os.Stderr, "moac: unknown command %v\n%s", cmd, usage)
		os.Exit(1)
	}
}
