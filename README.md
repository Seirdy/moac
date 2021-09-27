MOAC
====

[![godocs.io](https://godocs.io/git.sr.ht/~seirdy/moac?status.svg)](https://godocs.io/git.sr.ht/~seirdy/moac) [![builds.sr.ht status](https://builds.sr.ht/~seirdy/moac.svg)](https://builds.sr.ht/~seirdy/moac) [![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/5227/badge)](https://bestpractices.coreinfrastructure.org/projects/5227) [![Go Report Card](https://goreportcard.com/badge/codeberg.org/seirdy/moac)](https://goreportcard.com/report/codeberg.org/seirdy/moac)

[![sourcehut](https://img.shields.io/badge/repository-sourcehut-lightgrey.svg?logo=data:image/svg+xml;base64,PHN2ZyBmaWxsPSIjZmZmIiB2aWV3Qm94PSIwIDAgNTEyIDUxMiIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cGF0aCBkPSJNMjU2IDhDMTE5IDggOCAxMTkgOCAyNTZzMTExIDI0OCAyNDggMjQ4IDI0OC0xMTEgMjQ4LTI0OFMzOTMgOCAyNTYgOHptMCA0NDhjLTExMC41IDAtMjAwLTg5LjUtMjAwLTIwMFMxNDUuNSA1NiAyNTYgNTZzMjAwIDg5LjUgMjAwIDIwMC04OS41IDIwMC0yMDAgMjAweiIvPjwvc3ZnPg==)](https://sr.ht/~seirdy/MOAC) [![GitLab mirror](https://img.shields.io/badge/mirror-GitLab-orange.svg?logo=gitlab)](https://gitlab.com/Seirdy/moac) [![GitHub mirror](https://img.shields.io/badge/mirror-GitHub-black.svg?logo=github)](https://github.com/Seirdy/moac) [![Codeberg mirror](https://img.shields.io/badge/mirror-Codeberg-blue.svg?logo=codeberg)](https://codeberg.org/Seirdy/moac)

Generate passwords and measure their strength according to physical limits to computation.

This software is concerned only with password strength, and knows nothing about the context in which passwords will be used; as such, it makes the assumption that password guessability is the only metric that matters, and a brute-force attack is constrained only by the laws of physics. It's inspired by a blog post I wrote: [Becoming physically immune to brute-force attacks](https://seirdy.one/2021/01/12/password-strength.html).

Users provide given values like the mass available to attackers, a time limit for the brute-force attack, and the energy available. `moac` outputs the likelihood of a successful attack or the minimum password entropy for a possible brute-force failure. Entropy is calculated with the assumption that passwords are randomly generated.

`moac-pwgen` can also generate passwords capable of withstanding a brute-force attack limited by given physical quantities. Generated passwords are random and near impossible to memorize; these are suitable for non-human entry, e.g. auto-entry by a password manager.

My original intent when making this tool was to illustrate how easy it is to make a password whose strength is "overkill". It has since evolved into a generic password generator and evaluator.

Project Status
--------------

MOAC is actively developed as of September 2021.

Installation
------------

Check the [refs](https://git.sr.ht/~seirdy/moac/refs/) for changelogs; select a ref to download pre-built binaries. You can verify releases using the signing key [1E892DB2A5F84479](https://seirdy.one/publickey.asc) Alternatively, build MOAC from source:

### Build dependencies

- Go toolchain
- `make` (tested with GNU Make 4.x, `bmake`, and OpenBSD Make; should be fairly portable).
- `scdoc` (for building man pages)
- `git` (optional; for embedding the version number in the binary)

```sh
sudo make install-strip # Install in /usr/local/ by default; set PREFIX to change.
```

To upgrade or reinstall, run the above again; to uninstall, run `sudo make uninstall`

### Reproducible builds

Run `make dist-reprod` to build a tarball containing reproducible binaries. Binaries should be reproducible for a given Go toolchain; check `.builds/setup-toolchain.sh` to see which toolchain to use to reproduce binary artifacts built by the Fedora buildserver (except for the `linux-sanitizers` ones).

Usage (with three examples)
---------------------------

For full usage of the command-line executables, see the [`moac(1)`](https://git.sr.ht/~seirdy/moac/tree/master/item/doc/moac.1.scd) and [`moac-pwgen(1)`](https://git.sr.ht/~seirdy/moac/tree/master/item/doc/moac-pwgen.1.scd) man pages. Man page sources are in `doc/`.

### Bottlenecks and redundancy

If a value is provided _and_ that value can be computed from other given values, the computed value will replace the provided value if the computed value is a greater bottleneck.

If the user supplies both mass and energy, the given energy will be replaced with the mass-energy of the provided mass if the given mass-energy is lower.

If the user supplies both a password and a password entropy, the given entropy will be replaced with the calculated entropy of the provided password if the calculated entropy is lower. If the user does not supply entropy or the physical values necessary to calculate it, the default entropy is `256` (the key length of AES-256).

Time and energy are the two bottlenecks to computation; the final result will be based on whichever is a greater bottleneck. Unless the lower bound of the energy per guess is orders of magnitude below the Landauer limit or a non-default "power" is provided, energy should always be a greater bottleneck.

When physical quantities are not given, default physical quantities are the mass of the visible universe and the power required to achieve Bremermann's limit at the energy efficiency given by the Landauer limit.

### Example 1: a password the Earth cannot crack

The novel _The Hitchhiker's Guide to the Galaxy_ revealed the Earth to be a supercomputer built to understand "the answer to Life, the Universe, and Everything". The computation was supposed to finish sometime around now.

Let's assume this is a maximally efficient quantum computer powered by the Earth's mass-energy:

- Age of the Earth: ~4.6 billion years, or ~1.45e17 seconds
- Mass of the Earth: ~5.97e24 kg
- Temperature somewhere in the upper mantle: 1900 K

```console
$ moac -qm 5.97e24 -t 1.45e17 -T 1900 entropy-limit
408
```

Understanding the answer to Life, the Universe, and Everything requires less than `2^408` computations. If the same computer instead tried to brute-force a password, what kind of password might be out of its reach?

```console
$ moac-pwgen -qm 5.97e24 -t 1.45e17 -T 1900 ascii latin
ɮʠðʋsĳóʣ[5ȍìŒŞȨRɸÒ¨ůİȤ&ǒŘĥėǺʞĚʥ¼ɖƅ~ɛ\{ƸÝ4Ǎ6ő&Æ
```

If the same computer instead tried to guess the password `ɮʠðʋsĳóʣ[5ȍìŒŞȨRɸÒ¨ůİȤ&ǒŘĥėǺʞĚʥ¼ɖƅ~ɛ\{ƸÝ4Ǎ6ő&Æ`, there's a chance that it wouldn't have succeeded in time.

_Note: given that the Earth wasn't hollow during the book's opening, it's unlikely that the Earth consumed its own mass to compute. Further research is necessary: perhaps it used solar power, or secret shipments of tiny black-hole batteries? Organic life was supposed to provide a large part of its functionality, so maybe we should restrict ourselves to the Earth's biomass._

### Example 2: a computation powered by the universe

MOAC's default values are measured values for the whole observable universe: the observable universe's mass and the temperature of cosmic background radiation. If we burned all the mass in the universe, how many computations can we do?

```console
$ moac entropy-limit # classical computer
307
$ moac -q entropy-limit # quantum computer
615
```

A universal computer would be able to do the equivalent of cracking a password with 307 bits of entropy, or twice that if it were a quantum computer. What would such a password look like?

```console
$ moac-pwgen -s 307
\Go/I-MYTP#kTar>I.Oe)*^3)"yWh@/P<|v\@LBF ?-N{2/N
$ moac-pwgen -s 615
R,EZO.bZ6--wd HtQXw#Q,#&O#.Y.6*Y{&~*Wa|CAUZRO2TIbx>'k>vQ@GatFpu>,PNPZ.BUJ;uIm(JvXFOff-/6-JkA+^e
```

Be aware that current state-of-the-art encryption typically involves key-lengths of 256 bits, so anything more than that is probably a waste.

### Example 3: bad password rules

Some websites in this day and age still insist on [bad password rules](https://xkcd.com/936/), requiring and forbidding the use of certain sets of characters. Say a website requires a password between 8 and 16 characters; it needs one uppercase letter, one lowercase letter, and one ASCII symbol that isn't `$\!@`.

```console
$ moac-pwgen -L 16 uppercase lowercase "#%^&*()-_=+][}{';\":"
oi&YSt:FwdAO_y}^
```

Ideas for other programs that can use `moac`
--------------------------------------------

- A separate program to "benchmark" external password-generation programs/scripts by repeatedly running them and giving measurements of the worst output.
- A GUI
- Plugins for existing password managers. Account for key length used in encryption; if the key length is lower than the password entropy, the key length is the bottleneck.

FAQ
---

### Why did you make MOAC?

Two reasons: the blog post I wrote (linked at the top) got me itching to implement its ideas, and I also want to use a good password generator in a password manager I'm working on.

### How does MOAC measure password entropy?

It takes a very naive approach, assuming that any attacker is optimizing for randomly-generated passwords. More specifically, it measures password entropy as if `moac-pwgen` generated the password. All it does it guess which charsets are used and measure permutations of available characters for the given password length.

When generating passwords using custom charsets, it de-duplicates repeated characters and cross-charset overlap to avoid over-inflating password entropy estimates.

### Why do these passwords look impossible to memorize or type?

MOAC is not meant to be used to generate passwords to type by hand. It's intended to be used with a password manager that auto-enters or copies/pastes passwords for you.

For contexts in which you must enter a password manually (e.g. a full-disk encryption password entered during boot), use something else. I recommend a diceware-based passphrase in a language/charset compatible with your preferred available input method.

### Why are there so many weird characters in the generated passwords?

Those "weird characters" are configurable; check the man pages or GoDoc for more info. I admit that charsets like `ipaExtensions` were mostly added for fun, but they can be quite useful for detecting bugs in other software that accepts text input.

Starting with v0.3.2, password generation defaults to alphanumerics and basic QWERTY symbols. I figured that this is probably for the best, as long as most of us have to work with software that breaks when encountering non-QWERTY symbols. After all, everyone knows that password entry existed long before [languages besides English](https://blog.tdwright.co.uk/2018/11/06/anglocentrism-broke-my-tests-ignore-localisation-at-your-peril/) were invented.

### Why does MOAC's password generator try to use one character from each charset?

A lot of bad software mandates the usage of one character from a given charset (one number, one symbol, etc). This makes compliance with that bad software easier.

### Does MOAC support grapheme clusters?

When measuring password strength, MOAC counts code points (runes); it does not group grapheme clusters together. This is suitable for entropy calculations. Most password length requirements do not take grapheme clusters into account, so code points also make for a sufficient measure of password length.

The only times MOAC could ever encounter grapheme clusters are when given a password to analyze or when given a custom charset to use when generating passwords. In the former situation, it will count each code point in a cluster towards a password's length and custom charset size; in the latter situation, it will treat each code point as a distinct character.

The inability to work with grapheme clusters in custom charsets is a known limitation.

MOAC may may or may not learn to preserve grapheme clusters in the future. There are two reasons for the hesitation:

1. MOAC's password generation uses Unicode code points ("runes") as the fundamental building blocks of passwords
2. MOAC's main focus is password generation optimized for _non-human entry_ (e.g., auto-entry by a password manager); generating passwords containing grapheme clusters is generally not useful for such a use-case.

The use-case for a generic random-string generator certainly exists, but isn't a significant enough focus to demand this great a refactor.

### How does MOAC handle characters with special behavior, like non-printable or mark characters?

When reading custom charsets, the `moac-pwgen` program only accepts Unicode code points that correspond to printable characters, excluding marks (some software incorrectly estimates password length when presented with marks). Other code points are ignored and trigger a warning sent to STDERR. This ensures that generated passwords work well with fickle software and behave as predictably as possible.

`moac` and the public functions in the Go library do not restrict acceptable Unicode ranges; they should work with all Unicode code points.

### Why are MOAC's default values the values of the observable universe?

MOAC assumes that if a user doesn't provide a value, the value is as close to a non-bottleneck as it can possibly be. For instance, an attacker cannot have infinite mass energy: they must be limited by what's in the observable universe. Mass _is_ a bottleneck, even if not specified.

Caveat: the limits of the speed of light in the face of an expanding universe make it impossible to escape our local galactic cluster. Later, I might change the default values to those of our local cluster. I might also try to compute total energy to include sources besides matter.

Contributing
------------

Development is mostly happening in the v2 branch, which contains breaking changes.

See `CONTRIBUTING.md` in the repository root for contribution guidelines.

Alternatives
------------

- [libpwquality](https://github.com/libpwquality/libpwquality/)
- [zxcvbn](https://www.usenix.org/conference/usenixsecurity16/technical-sessions/presentation/wheeler)
- [pwgen](http://sf.net/projects/pwgen)
- [cracklib](https://github.com/cracklib/cracklib)
- The password generator/evaluator in [KeePassXC](https://keepassxc.org/)
- [go-password-validator](https://github.com/wagslane/go-password-validator)
- [go-password](https://github.com/sethvargo/go-password)

