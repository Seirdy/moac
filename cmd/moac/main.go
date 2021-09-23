package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"

	"git.sr.ht/~seirdy/moac"
	"git.sr.ht/~seirdy/moac/entropy"
	"git.sr.ht/~seirdy/moac/internal/cli"
	"git.sr.ht/~seirdy/moac/internal/version"
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

func getBruteForceability(givens *moac.Givens, quantum bool) float64 {
	likelihood, err := moac.BruteForceability(givens, quantum)
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

	computedEntropy, err := entropy.Entropy(givens.Password)
	cli.ExitOnErr(err, "")

	return computedEntropy
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
	minEntropy, err := moac.MinEntropy(givens, quantum)
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
	if len(args) == 0 {
		fmt.Printf("%.3g\n", getBruteForceability(givens, quantum))
		os.Exit(0)
	}

	cmd := args[0]
	switch cmd {
	case "strength":
		fmt.Printf("%.3g\n", getBruteForceability(givens, quantum))
	case "entropy":
		fmt.Printf("%.3g\n", getEntropy(givens))
	case "entropy-limit":
		fmt.Printf("%.3g\n", getMinEntropy(givens, quantum))
	default:
		fmt.Fprintf(os.Stderr, "moac: unknown command %v\n%s", cmd, usage)
		os.Exit(1)
	}
}
