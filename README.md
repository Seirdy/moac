MOAC
====

[![godocs.io](https://godocs.io/git.sr.ht/~seirdy/moac?status.svg)](https://godocs.io/git.sr.ht/~seirdy/moac) [![builds.sr.ht status](https://builds.sr.ht/~seirdy/moac.svg)](https://builds.sr.ht/~seirdy/moac)

[![sourcehut](https://img.shields.io/badge/repository-sourcehut-lightgrey.svg?logo=data:image/svg+xml;base64,PHN2ZyBmaWxsPSIjZmZmIiB2aWV3Qm94PSIwIDAgNTEyIDUxMiIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cGF0aCBkPSJNMjU2IDhDMTE5IDggOCAxMTkgOCAyNTZzMTExIDI0OCAyNDggMjQ4IDI0OC0xMTEgMjQ4LTI0OFMzOTMgOCAyNTYgOHptMCA0NDhjLTExMC41IDAtMjAwLTg5LjUtMjAwLTIwMFMxNDUuNSA1NiAyNTYgNTZzMjAwIDg5LjUgMjAwIDIwMC04OS41IDIwMC0yMDAgMjAweiIvPjwvc3ZnPg==)](https://sr.ht/~seirdy/MOAC) [![GitLab mirror](https://img.shields.io/badge/mirror-GitLab-orange.svg?logo=gitlab)](https://gitlab.com/Seirdy/moac) [![GitHub mirror](https://img.shields.io/badge/mirror-GitHub-black.svg?logo=github)](https://github.com/Seirdy/moac) [![Codeberg mirror](https://img.shields.io/badge/mirror-Codeberg-blue.svg?logo=codeberg)](https://codeberg.org/Seirdy/moac)

Generate passwords and measure their strength according to physical limits to computation.

This software is concerned only with password strength, and knows nothing about the context in which passwords will be used; as such, it makes the assumption that password guessability is the only metric that matters, and a brute-force attack is constrained only by the laws of physics. It's inspired by a blog post I wrote: [Becoming physically immune to brute-force attacks](https://seirdy.one/2021/01/12/password-strength.html).

Users provide given values like the mass available to attackers, a time limit for the brute-force attack, and the energy available. `moac` outputs the likelihood of a successful attack or the minimum password entropy for a possible brute-force failure. Entropy is calculated with the assumption that passwords are randomly generated.

`moac-pwgen` can also generate passwords capable of withstanding a brute-force attack limited by given physical quantities.

My original intent when making this tool was to illustrate how easy it is to make a password whose strength is "overkill". It has since evolved into a generic password generator and evaluator.

Project Status
--------------

MOAC is actively developed as of September 2021. It's almost ready for a v1.0.0 release; I'd just like to hear some opinions/feedback before I tag+push.

Installation
------------

### Dependencies

- Go toolchain
- `make` (tested with GNU Make, `bmake`, and OpenBSD Make)
- `scdoc` (for building manpages)

```sh
sudo make install # Install in /usr/local/ by default
```

Usage (with examples)
---------------------

For full usage of the command-line executables, see the [`moac(1)`](https://git.sr.ht/~seirdy/moac/tree/master/item/doc/moac.1.scd) and [`moac-pwgen(1)`](https://git.sr.ht/~seirdy/moac/tree/master/item/doc/moac-pwgen.1.scd) manpages. Manpage sources are in `doc/`.

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

_Note: given that the Earth wasn't hollow during the book's opening, it's unlikely that the Earth consumed its own mass to compute. The Earth is also much warmer than 2.7 kelvins. Further research is necessary: perhaps it used solar power, or secret shipments of tiny black-hole batteries? Organic life was supposed to provide a large part of its functionality, so maybe we should restrict ourselves to the Earth's biomass._

### Ideas for other programs that can use `moac`

- A separate program to "benchmark" external password-generation programs/scripts by repeatedly running them and giving measurements of the worst output.
- A GUI
- Plugins for existing password managers. Account for key length used in encryption; if the key length is lower than the password entropy, the key length is the bottleneck.

FAQ
---

### Why did you make MOAC?

Two reasons: the blog post I wrote (linked at the top) got me itching to implement its ideas, and I also want to use a good password generator in a password manager I'm working on.

### How does MOAC measure password entropy?

It takes a very naive approach, assuming that any attacker is optimizing for randomly-generated passwords. More specifically, it measures password entropy as if `moac-pwgen` generated the password. All it does it guess which charsets are used and measure permutations of available characters for the given password length.

### Why do these passwords look impossible to memorize or type?

MOAC is not meant to be used to generate passwords to type by hand. It's intended to be used with a password manager that auto-types or copies passwords for you.

For contexts in which you can't paste a password (e.g. a full-disk encryption password entered during boot), use something else.

### Why are there so many weird characters in the generated passwords?

Those "weird characters" are configurable; check the manpages or GoDoc for more info. I admit that charsets like `ipaExtensions` were mostly added for fun, but they can be quite useful for detecting bugs in other software that accepts text input.

Starting with v0.3.2, password generation defaults to alphanumerics and basic QWERTY symbols. I figured that this is probably for the best, as long as most of us have to work with software that breaks when encountering non-QWERTY symbols. After all, everyone knows that password entry existed long before [languages besides English](https://blog.tdwright.co.uk/2018/11/06/anglocentrism-broke-my-tests-ignore-localisation-at-your-peril/) were invented.

Contributing
------------

I try to accommodate everyone's workflow. Here's ways to contribute and how, in my order of preference:

### Bug reports

Preferred and canonical location: <https://todo.sr.ht/~seirdy/MOAC>. Send an email to <mailto:~seirdy/MOAC@todo.sr.ht> to automatically file a bug, no account needed.

I also check issues in the GitHub, GitLab, and Codeberg mirrors linked at the top of the README, if you prefer. No matter which option you choose, your bug gets emailed to me.

### Patches, questions, and feature requests

Preferred location: <https://lists.sr.ht/~seirdy/moac>. Send emails and patches to <mailto:~seirdy/moac@lists.sr.ht>.

For patches, run `make test` to run the unit tests, and `make lint` to run [golangci-lint](https://github.com/golangci/golangci-lint) and [gokart](https://github.com/praetorian-inc/gokart).

I also check the GitHub, GitLab, and Codeberg mirrors for issues and PRs.

### Direct contact

Check my contact info at the bottom of the manpages.

Alternatives
------------

- [libpwquality](https://github.com/libpwquality/libpwquality/)
- [zxcvbn](https://www.usenix.org/conference/usenixsecurity16/technical-sessions/presentation/wheeler)
- [pwgen](http://sf.net/projects/pwgen)
- [cracklib](https://github.com/cracklib/cracklib)
- The password generator/evaluator in [KeePassXC](https://keepassxc.org/)

