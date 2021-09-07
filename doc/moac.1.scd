moac(1)

# NAME

moac - analyze password strength with physical limits

# SYNOPSIS

moac [OPTIONS...] [COMMAND]

# OPTIONS

*-h*
	Display a brief help message.

*-q*
	Account for quantum computers using Grover's algorithm

*-r*
	Interactively enter a password in the terminal; overrides *-p*

*-e* <energy>
	Maximum energy used by attacker (J).

*-s* <entropy>
	Password entropy. Overrides the provided password's computed entropy if it's
	lower.

*-m* <mass>
	Mass at attacker's disposal (kg). Used to compute mass-energy. Overrides the
	value of *-e* if the computed mass-energy is lower.

*-g* <energy>
	Energy used per guess (J).

*-P* <power>
	Power available to the computer (W)

*-t* <time>
	Time limit for brute-force attack (s).

*-p* <password>
	Password to analyze. Note that this will be visible in plaintext to anything
	on your system that can view running processes, so don't use a real password.
	Use *-r* to interactively enter a password instead.

# DEFAULT VALUES

When a physical value isn't given and cannot be computed from other given
values, *moac* uses the following default values:

*entropy*: 256

*mass*: mass of the visible universe

*energy per guess*: Landauer limit

*guesses per second*: product of Bremermann's Limit and mass.

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
