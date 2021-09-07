MOAC
====

[![godocs.io](https://godocs.io/git.sr.ht/~seirdy/moac?status.svg)](https://godocs.io/git.sr.ht/~seirdy/moac)

[![sourcehut](https://img.shields.io/badge/repository-sourcehut-lightgrey.svg?logo=data:image/svg+xml;base64,PHN2ZyBmaWxsPSIjZmZmIiB2aWV3Qm94PSIwIDAgNTEyIDUxMiIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cGF0aCBkPSJNMjU2IDhDMTE5IDggOCAxMTkgOCAyNTZzMTExIDI0OCAyNDggMjQ4IDI0OC0xMTEgMjQ4LTI0OFMzOTMgOCAyNTYgOHptMCA0NDhjLTExMC41IDAtMjAwLTg5LjUtMjAwLTIwMFMxNDUuNSA1NiAyNTYgNTZzMjAwIDg5LjUgMjAwIDIwMC04OS41IDIwMC0yMDAgMjAweiIvPjwvc3ZnPg==)](https://sr.ht/~seirdy/MOAC) [![GitLab mirror](https://img.shields.io/badge/mirror-GitLab-orange.svg?logo=gitlab)](https://gitlab.com/Seirdy/moac) [![GitHub mirror](https://img.shields.io/badge/mirror-GitHub-black.svg?logo=github)](https://github.com/Seirdy/moac) [![Codeberg mirror](https://img.shields.io/badge/mirror-Codeberg-blue.svg?logo=codeberg)](https://codeberg.org/Seirdy/moac)

[![builds.sr.ht status](https://builds.sr.ht/~seirdy/moac.svg)](https://builds.sr.ht/~seirdy/moac)

`moac` is a tool that takes a unique approach to generating passwords and analyzing their strength. It's concerned only with password strength, and knows nothing about the context in which passwords will be used; as such, it makes the assumption that password guessability is the only metric that matters, and a brute-force attack is constrained only by the laws of physics. It's inspired by a blog post I wrote: [Becoming physically immune to brute-force attacks](https://seirdy.one/2021/01/12/password-strength.html).

Users provide given values like the mass available to attackers, a time limit for the brute-force attack, and the energy available. `moac` outputs the likelihood of a successful attack or the minimum password entropy for a possible brute-force failure. Entropy is calculated with the assumption that passwords are randomly generated.

`moac` can also generate passwords capable of withstanding a brute-force attack limited by given physical quantities.

My original intent when making this tool was to illustrate how easy it is to make a password whose strength is "overkill". It has since evolved into a generic password generator and evaluator.

**Note: until version 1.0.0 is released, MOAC is only suitable for educational/exploratory use and should not be considered stable. Do not use it with your actual passwords yet.**

Installation
------------

```sh
make install
```

Usage
-----

For full usage of the command-line executables, see `moac(1)` and `moac-pwgen(1)`. Manpages are in `doc/`.

### Bottlenecks and redundancy

If a value is provided _and_ that value can be computed from other given values, the computed value will replace the provided value if the computed value is a greater bottleneck.

If the user supplies both mass and energy, the given energy will be replaced with the mass-energy of the provided mass if the given mass-energy is lower.

If the user supplies both a password and a password entropy, the given entropy will be replaced with the calculated entropy of the provided password if the calculated entropy is lower. If the user does not supply entropy or the physical values necessary to calculate it, the default entropy is `256` (the key length of AES-256).

Time and energy are the two bottlenecks to computation; the final result will be based on whichever is a greater bottleneck. Unless the lower bound of the energy per guess is orders of magnitude below the Landauer limit, energy should always be a greater bottleneck.

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
$ moac-pwgen -qm 5.97e24 -t 1.45e17 lowercase uppercase numbers symbols latin
ɥìƄ¦sČÍM²ȬïľA\ɻ¨zŴǓĤúǓ¤ʬƗ;ɮĢƃƅǞɃƜʌȴɖǃƨǥ_Ǝ3ſǹǅɃ8ɟ
```

If the same computer instead tried to guess the password `,ȿĢıqɽȂīĲďɖȟMǧiœcɪʊȦĻțșŌƺȰ&ǡśŗȁĵɍɞƋIŀƷ?}ʯ4ůʑʅęȳŞ`, there's a chance that it wouldn't have succeeded in time.

_Note: given that the Earth wasn't hollow during the book's opening, it's unlikely that the Earth consumed its own mass to compute. Further research is necessary; perhaps it used solar power, or secret shipments of tiny black-hole batteries? Organic life was supposed to provide a large part of its functionality, so maybe we should restrict ourselves to the Earth's biomass._

Roadmap
-------

### Roadmap for 1.0.0

The actual code:

- [X] More comprehensive tests: cover everything that should be reachable
- [X] Move password generation to its own sub-package
- [X] CLI: ~~Separate global and command-specific options~~ split pwgen into own executable
- [X] Library: API seems finalized for 1.0

Other stuff:

- [X] CI/CD
- [X] Manpage for CLI
- [ ] Shell completion
- [ ] Set up signed releases or distro packages

Last steps before releasing v1.0.0:

- [ ] Get `moac`'s code reviewed by some people with more experience in software security.
- [ ] Link to it in my old blog post on brute-force immunity

### Future

- [ ] Account for quantum memory changing the constraints on energy-bound computation speed
- [ ] Estimate computations per guess, possibly using additional context like KDFs used

### Ideas for other programs that can use `moac`

- A separate program to "benchmark" external password-generation programs/scripts by repeatedly running them and giving measurements of the worst output.
- A GUI
- Plugins for existing password managers. Account for key length used in encryption; if the key length is lower than the password entropy, the key length is the bottleneck.

Alternatives
------------

- [libpwquality](https://github.com/libpwquality/libpwquality/)
- [zxcvbn](https://www.usenix.org/conference/usenixsecurity16/technical-sessions/presentation/wheeler)
- [pwgen](http://sf.net/projects/pwgen)
- [cracklib](https://github.com/cracklib/cracklib)
- The password generator/evaluator in [KeePassXC](https://keepassxc.org/)

License
-------

Copyright (C) 2021 Rohan Kumar

This program is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

