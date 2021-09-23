moac(1)

# NAME

moac - analyze password strength with physical limits

# SYNOPSIS

moac [OPTIONS...] [COMMAND]

# OPTIONS

*-h*
	Display a brief help message and exit.

*-v*
	Print the version and exit.

*-q*
	Account for quantum computers using Grover's algorithm. Note that as of late
	2021, the viability of quantum computers capable of using Grover's Algorithm
	is unknown.

*-r*
	Interactively enter a password in the terminal; overrides *-p*

*-e* <energy>
	Maximum energy used by attacker (J).

*-s* <entropy>
	Password entropy. Overrides the provided password's computed entropy if it's
	lower. Defaults to 256 bits.

*-m* <mass>
	Mass at attacker's disposal (kg). Used to compute mass-energy. Overrides the
	value of *-e* if the computed mass-energy is lower.

*-g* <energy>
	Energy used per guess (J).

*-P* <power>
	Power available to the computer (W)

*-T* <temperature>
	The temperature of the system (K). Defaults to 273 kelvin. Used to calculate
	the Landauer limit.

	Temperature must be between 0 and the Planck Temperature, noninclusive.

*-t* <time>
	Time limit for brute-force attack (s).

*-p* <password>
	Password to analyze. Use '-' to read the password from stdin.

	Note that when passing a password as an argument, it will be visible in
	plaintext to anything on your system that can view running processes. Use *-r*
	to interactively enter a password instead.

# DEFAULT VALUES

When a physical value isn't given and cannot be computed from other given
values, *moac* uses the following default values:

*entropy*: 256

*mass*: mass of the visible universe.

*energy per guess*: Landauer limit.

*guesses per second*: product of Bremermann's Limit and mass.

*temperature*: 2.7 K, a low estimate for the temperature of cosmic background
radiation.

# COMMANDS

*strength*
	Calculate the liklihood of a successful guess.

*entropy*
	Calculate the entropy of the given password

*entropy-limit*
	Calculate the minimum entropy for a brute-force attack failure.

If no command is given, default to *strength*.

# SEE ALSO

*moac-pwgen*(1)

# AUTHOR

Maintained by Rohan Kumar <seirdy@seirdy.one>. For more information about moac
development, see https://sr.ht/~seirdy/MOAC
