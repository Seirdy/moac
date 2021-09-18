moac-pwgen(1)

# NAME

moac-pwgen - generate passwords with the described strength

# SYNOPSIS

moac-pwgen [OPTIONS...] [CHARSETS...]

# OPTIONS

*-h*
	Display a brief help message and exit.

*-v*
	Print the version and exit.

*-q*
	Account for quantum computers using Grover's algorithm

*-e* <energy>
	Maximum energy used by attacker (J).

*-s* <entropy>
	Desired assword entropy. Defaults to 256 bits.

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

*-t* <time>
	Time limit for brute-force attack (s).

*-l* <length>
	Minimum number of characters in generated password. Overrides value of *-s* if
	doing so would increase password length.

*-L* <length>
	Maximum number of characters in generated password. Overrides value of *-s* if
	doing so would decrease password length.

# DEFAULT VALUES

See the *DEFAULT VALUES* section in *moac*(1) for physical values.

The default charset is *ascii*.

# CHARSETS

Generated passwords contain at least one character from each charset specified.
Choose from the following charsets:

*lowercase*	[a-z]

*uppercase*	[A-Z]

*numbers*	[0-9]

*symbols*	!"#%&'()\*+,-./:;<=>?@[\\]^\_`{|}~$-

*latin1*	the "Latin-1" Unicode block

*latinExtendedA*	the "Latin Extended A" Unicode block

*latinExtendedB*	the "Latin Extended B" Unicode block

*ipaExtensions*	the "IPA Extensions" Unicode block

*ascii*	alias for "lowercase uppercase numbers symbols". This is the default.

*latin*	alias for "latin1 latinExtendedA latinExtendedB ipaExtensions"

Any argument after *[OPTIONS]* that is not listed above is assumed to be a
string containing characters of a charset.

# SEE ALSO

*moac*(1)

# AUTHOR

Maintained by Rohan Kumar <seirdy@seirdy.one>. For more information about moac
development, see https://sr.ht/~seirdy/MOAC
