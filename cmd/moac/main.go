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
	helpText = "moac - analyze password strength with physical limits\n" + usage
)

func parseOpts( //nolint:cyclop // complexity solely determined by cli flag count
	opts *[]getopt.Option,
) (givens moac.Givens, quantum, readPassword bool, err error) {
	for _, opt := range *opts {
		switch opt.Option {
		case 'h':
			fmt.Fprint(os.Stderr, helpText)
			os.Exit(0)
		case 'v':
			fmt.Println(cli.GetVersion())
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
			err = fmt.Errorf("%w: invalid value for -%c: %s\n%s", cli.ErrBadCmdline, opt.Option, opt.Value, usage)

			break
		}
	}

	return givens, quantum, readPassword, err
}

func readPwInteractive() (pw string, err error) {
	fmt.Print("Enter password: ")

	bytepw, err := term.ReadPassword(int(syscall.Stdin)) //nolint:unconvert // needed for some platforms

	fmt.Println()

	if err != nil {
		return pw, fmt.Errorf("failed to read pw interactively: %w", err)
	}

	pw = string(bytepw)

	return pw, nil
}

func readPwStdin() (pw string, err error) {
	stdinBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return pw, fmt.Errorf("faild to read pw from stdin: %w", err)
	}

	return string(stdinBytes), nil
}

func main() {
	os.Exit(main1())
}

func main1() int {
	output, err := getOutput()
	if cli.DisplayErr(err, "") {
		fmt.Printf(cli.FloatFmt, output)

		return 0
	}

	return 1
}

func getOutput() (output float64, err error) {
	opts, optind, err := getopt.Getopts(os.Args, "hvqre:s:m:g:P:T:t:p:")
	if err != nil {
		return output, fmt.Errorf("%w\n%s", err, usage)
	}

	givens, quantum, readPassword, err := parseOpts(&opts)
	if err != nil {
		return output, fmt.Errorf("moac: %w", err)
	}

	if readPassword {
		givens.Password, err = readPwInteractive()
	} else if givens.Password == "-" {
		givens.Password, err = readPwStdin()
	}

	if err != nil {
		return output, err
	}

	givens.Password = strings.TrimSuffix(givens.Password, "\n")
	args := os.Args[optind:]
	cmd := "strength"

	if len(args) > 0 {
		cmd = args[0]
	}

	return runCmd(cmd, quantum, &givens)
}

func runCmd(cmd string, quantum bool, givens *moac.Givens) (output float64, err error) {
	switch cmd {
	case "strength":
		if quantum {
			output, err = givens.BruteForceabilityQuantum()
		} else {
			output, err = givens.BruteForceability()
		}
	case "entropy":
		if givens.Password == "" {
			err = fmt.Errorf("%w: missing password", moac.ErrMissingValue)
		}

		output = entropy.Entropy(givens.Password)
	case "entropy-limit":
		if quantum {
			output, err = givens.MinEntropyQuantum()
		} else {
			output, err = givens.MinEntropy()
		}
	default:
		err = fmt.Errorf("%w: unknown command %v", cli.ErrBadCmdline, cmd)
	}

	return output, err
}
