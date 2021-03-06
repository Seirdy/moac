moac
====

[![godocs.io](https://godocs.io/git.sr.ht/~seirdy/moac?status.svg)](https://godocs.io/git.sr.ht/~seirdy/moac) [![sourcehut](https://img.shields.io/badge/repository-sourcehut-lightgrey.svg?logo=data:image/svg+xml;base64,PHN2ZyBmaWxsPSIjZmZmIiB2aWV3Qm94PSIwIDAgNTEyIDUxMiIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cGF0aCBkPSJNMjU2IDhDMTE5IDggOCAxMTkgOCAyNTZzMTExIDI0OCAyNDggMjQ4IDI0OC0xMTEgMjQ4LTI0OFMzOTMgOCAyNTYgOHptMCA0NDhjLTExMC41IDAtMjAwLTg5LjUtMjAwLTIwMFMxNDUuNSA1NiAyNTYgNTZzMjAwIDg5LjUgMjAwIDIwMC04OS41IDIwMC0yMDAgMjAweiIvPjwvc3ZnPg==)](https://sr.ht/~seirdy/MOAC) [![GitLab mirror](https://img.shields.io/badge/mirror-GitLab-orange.svg?logo=gitlab)](https://gitlab.com/Seirdy/moac) [![GitHub mirror](https://img.shields.io/badge/mirror-GitHub-black.svg?logo=github)](https://github.com/Seirdy/moac)

`moac` is a tool to analyze password strength given physical limits to computation. It's inspired by a blog post I wrote: [Becoming physically immune to brute-force attacks](https://seirdy.one/2021/01/12/password-strength.html).

Users provide given values like the mass available to attackers, a time limit for the brute-force attack, and the energy available. `moac` outputs the likelihood of a successful attack or the minimum password entropy for a possible brute-force failure.

`moac` can also generate passwords capable of withstanding a brute-force attack limited by given physical quantities.

[zxcvbn-go](https://github.com/nbutton23/zxcvbn-go) calculates password entropy.

Installation
------------

```sh
GO111MODULE=on go install git.sr.ht/~seirdy/moac/cmd/moac@latest
```

Usage
-----

```
moac - analyze password strength with physical limits

USAGE:
  moac [OPTIONS] [COMMAND] [ARGS]

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
```

### Bottlenecks and redundancy

If a value is provided _and_ that value can be computed from other given values, the computed value will replace the provided value if the computed value is a greater bottleneck.

If the user supplies both mass and energy, the given energy will be replaced with the mass-energy of the provided mass if the given mass-energy is lower.

If the user supplies both a password and a password entropy, the given entropy will be replaced with the calculated entropy of the provided password if the calculated entropy is lower.

Time and energy are the two bottlenecks to computation; the final result will be based on whichever is a greater bottleneck. With the default energy per guess (the Landauer limit), energy should always be a greater bottleneck.

When physical quantities are not given, default physical quantities are the mass of the visible universe and the power required to achieve Bremermann's limit at the energy efficiency given by the Landauer limit.

### Example: a password the Earth cannot crack

The novel _The Hitchhiker's Guide to the Galaxy_ revealed the Earth to be a supercomputer built to understand "the answer to Life, the Universe, and Everything". The computation was supposed to finish sometime around now.

Let's assume this is a maximally efficient quantum computer powered by the Earth's mass-energy:

- Age of the Earth: ~4.6 billion years, or ~1.45e17 seconds
- Mass of the Earth: ~5.97e24 kg

```console
$ moac -qm 5.97e24 -t 1.45e17 entropy-limit
427
```

Understanding the answer to Life, the Universe, and Everything requires less than `2^427` computations. If the same computer instead tried to brute-force a password, what kind of password might be out of its reach?

```console
$ moac -qm 5.97e24 -t 1.45e17 pwgen lowercase uppercase numbers symbols extendedASCII
v¢JÊÙúQ§4mÀÛªZûYÍé©mËiÐ× "½J6y.ñíí'è¦ïÏµ°~
```

If the same computer instead tried to guess the password `v¢JÊÙúQ§4mÀÛªZûYÍé©mËiÐ× "½J6y.ñíí'è¦ïÏµ°~`, there's a chance that it wouldn't have succeeded in time.

_Note: given that the Earth wasn't hollow during the book's opening, it's unlikely that the Earth consumed its own mass to compute. Further research is necessary; perhaps it used solar power, or secret shipments of tiny black-hole batteries? Organic life was supposed to provide a large part of its functionality, so maybe we should restrict ourselves to the Earth's biomass._

### Example: skip the physics generate/measure passwords

Roadmap for 0.1.0
-----------------

Consider this project highly unstable before v0.1.0.

- godocs and manpage
- Better error handling: validate input, etc.
- Write tests.
- Analyze custom charsets to see if any characters fit in known charsets

Roadmap for 0.2.0
-----------------

- Read from a config file.
- Add a command to output requirements for a brute-force attack (time/energy/mass required) with the given constraints.
- zxcvbn-go has a lot of functionality that `moac` doesn't need; write an entropy estimator that's a bit simpler but gives similar results, optimized for pseudorandom passwords (no dictionary words, focus on estimating charset size and repetitions/patterns).

### Ideas for other programs that can use `moac`

- A separate program to "benchmark" external password-generation programs/scripts by repeatedly running them and giving measurements of the worst output.
- A GUI
- Plugins for existing password managers. Account for key length used in encryption; if the key length is lower than the password entropy, the key length is the bottleneck.

License
-------

Copyright (C) 2021 Rohan Kumar

This program is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

