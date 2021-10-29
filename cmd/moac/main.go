package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"

	"git.sr.ht/~seirdy/moac/v2"
	"git.sr.ht/~seirdy/moac/v2/entropy"
	"git.sr.ht/~seirdy/moac/v2/internal/cli"
	"git.sr.ht/~sircmpwn/getopt"
	"golang.org/x/term"
	"github.com/rivo/uniseg"
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
  -G <guesses>	Guesses per second in a brute-force attack
  -p <password>	Password to analyze; use "-" for stdin

COMMANDS:
  strength	Calculate the liklihood of a successful guess 
  entropy	Calculate the entropy of the given password
  entropy-limit	Calculate the minimum entropy for a brute-force attack failure.
`
	helpText = "moac - analyze password strength with physical limits\n" + usage
)

// Any changes to the below const values would be a breaking change.
const (
	strengthCmd   = "strength"
	entropyCmd    = "entropy"
	minEntropyCmd = "entropy-limit"
)

func parseOpts(
	opts *[]getopt.Option,
) (givens moac.Givens, quantum, exitEarly bool, err error) {
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
		case 'r':
			givens.Password, err = readPwInteractive()
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
		case 'p':
			givens.Password = opt.Value
		}

		if err != nil {
			err = fmt.Errorf("%w: invalid value for -%c: %s\n%s", cli.ErrBadCmdline, opt.Option, opt.Value, usage)

			break
		}

		if exitEarly {
			break
		}
	}

	return givens, quantum, exitEarly, err
}

func readPwInteractive() (pw string, err error) {
	fmt.Print("Enter password: ")

	bytepw, err := term.ReadPassword(int(syscall.Stdin)) //nolint:unconvert // needed for some platforms

	fmt.Println()

	if err != nil {
		return pw, fmt.Errorf("failed to read pw interactively: %w", err)
	}

	return string(bytepw), nil
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
	output, exitEarly, err := getOutput()
	if exitEarly {
		return 0
	}

	if cli.DisplayErr(err, "") {
		fmt.Printf(cli.FloatFmt, output)

		return 0
	}

	return 1
}

func getOutput() (output float64, exitEarly bool, err error) {
	opts, optind, err := getopt.Getopts(os.Args, "hvqre:s:m:g:P:T:t:G:p:")
	if err != nil {
		return output, exitEarly, fmt.Errorf("%w\n%s", err, usage)
	}

	givens, quantum, exitEarly, err := parseOpts(&opts)
	if err != nil || exitEarly {
		return output, exitEarly, err
	}

	givens.Password, err = processPassword(givens.Password)
	if err != nil {
		return output, exitEarly, fmt.Errorf("moac: %w", err)
	}

	graphemes := uniseg.NewGraphemes(givens.Password)
	for graphemes.Next() {
		if len(graphemes.Runes()) > 1 {
			fmt.Fprintf(os.Stderr, "warning: charsets contain grapheme clusters, will be treated as distinct codepoints\n")
			break
		}
	}

	cmd := strengthCmd

	if len(os.Args) > optind {
		cmd = os.Args[optind]
	}

	switch {
	case cmd == entropyCmd: // entropy is independent of Grover search feasibility.
		output, err = getEntropy(givens.Password)
	case quantum:
		output, err = runCmdQuantum(cmd, &givens)
	default:
		output, err = runCmdClassical(cmd, &givens)
	}

	return output, exitEarly, err
}

func processPassword(oldPw string) (newPw string, err error) {
	newPw = oldPw
	if newPw == "-" {
		newPw, err = readPwStdin()
	}

	for len(newPw) > 0 && newPw[len(newPw)-1] == '\n' {
		newPw = newPw[:len(newPw)-1]
	}

	return newPw, err
}

func runCmdClassical( //nolint:dupl // this duplication is worth keeping flat switches imo
	cmd string, givens *moac.Givens) (output float64, err error) {
	switch cmd {
	case strengthCmd:
		output, err = givens.BruteForceability()
	case minEntropyCmd:
		output, err = givens.MinEntropy()
	default:
		err = fmt.Errorf("%w: unknown command %v", cli.ErrBadCmdline, cmd)
	}

	return output, err
}

func runCmdQuantum( //nolint:dupl // see above
	cmd string, givens *moac.Givens) (output float64, err error) {
	switch cmd {
	case strengthCmd:
		output, err = givens.BruteForceabilityQuantum()
	case minEntropyCmd:
		output, err = givens.MinEntropyQuantum()
	default:
		err = fmt.Errorf("%w: unknown command %v", cli.ErrBadCmdline, cmd)
	}

	return output, err
}

func getEntropy(password string) (res float64, err error) {
	if password == "" {
		err = fmt.Errorf("%w: missing password", moac.ErrMissingValue)
	}

	return entropy.Entropy(password), err
}
