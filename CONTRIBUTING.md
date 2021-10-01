Contributing to MOAC
====================

Ways to contribute
------------------

I try to accommodate everyone's workflow. Here's ways to contribute and how, in my order of preference:

### Security vulnerabilities, non-public contact

If you want to connect with me directly: my email address is at the bottom of the man pages or in the commit logs. Vulnerability disclosures should be PGP-encrypted using the PGP key [1E892DB2A5F84479](https://seirdy.one/publickey.asc).

Alternatively, send an encrypted message on Matrix to `@seirdy:envs.net`

### Bug reports and TODOs

Preferred and canonical ticket tracker: <https://todo.sr.ht/~seirdy/MOAC>. Send an email to [~seirdy/moac@todo.sr.ht](mailto:~seirdy/MOAC@todo.sr.ht) to automatically file a bug, no account needed. The tracker might also have some tickets labeled "good first issue" ideal for contributors with less experience.

I also check issues in the GitHub, GitLab, and Codeberg mirrors linked at the top of the README, if you prefer. No matter which option you choose, your bug gets emailed to me.

### Patches, questions, and feature requests

Preferred location: <https://lists.sr.ht/~seirdy/moac>. Send emails and patches to [~seirdy/moac@lists.sr.ht](mailto:~seirdy/moac@lists.sr.ht). I also check the GitHub, GitLab, and Codeberg mirrors for issues and PRs.

#### Coding standards

Contributions don't need to follow these standards to be useful. If a useful patch doesn't pass the below checks, I might clean it up myself.

This project uses `gofumpt`, `fieldalignment`, `shfmt`, and `mdfmt -stxHeaders` for formatting; it uses `golangci-lint`, `gokart`, `go-consistent`, `shfmt` (again), and `checkmake` for linting. You can install all of them to your `GOBIN` by running `.builds/install-linters.sh`

Run `make fmt` to format code, `make lint` to run the linters (except `mdfmt`), and `make test` to run unit tests. `make pre-commit` runs all three. I recommend using [committer](https://github.com/Gusto/committer) to auto-run pre-commit checks; just add `committer` to your hooks.

The linters are very opinionated. If you find this annoying, you can send your patch anyway; I'll clean it up if it looks useful.

See the "Testing" section near the bottom for info about the tests.

### Other ways to help

- See if you can/can't reproduce binaries for a given installation of the Go toolchain, and share your findings to the mailing list or GitHub/GitLab/Codeberg issue trackers.
- Check out the documentation and see if anything seems unclear; I lack the perspective of someone reading these docs for the first time.

Quick architecture overview
---------------------------

Excluding tests, MOAC has <1k SLOC; it shouldn't be hard to grok. Here's a one-minute overview:

- `givens.go` handles given physical values (what you'd call "the givens" if you were solving a physics problem) and computes missing values/bottlenecks.
- `charsets` handles parsing, building, and de-duplicating charsets to use when calculating password entropy or building passwords.
- `entropy`, well, calculates entropy. It figures out what charsets are contained in a password (saving these in a data structure defined by `charsets`) and figures out how many combinations can fit in the resulting space.
- `pwgen` contains the `GenPW` function builds passwords that match the given requirements: length bounds, target entropy, and charsets to use.

Testing
-------

For the library: everything possible should be covered by tests. If a branch that handles an error should be impossible to reach and is therefore uncovered, replace it with a panic to indicate the presence of a bug. Any uncovered line that isn't a panic is in need of a test.

For the CLI: this uses [testscript](https://godocs.io/github.com/rogpeppe/go-internal/testscript) to test CLI behavior.

If you want live test feedback while hacking and find the tests to be too slow (they typically take under 3s by default on my low-end notebook), set the environment variable `LOOPS` to something below `64`; running `make test-quick` will set it to `10`. Test-cases for password generation run multiple times because of the non-determinism inherent to random password generation. Tests are a bit slow since `GenPW()`'s tests have thousands of test-cases generated from combinations of possible parameters.

If you notice that a change causes a big slowdown in `make test`, run `make test-prof` to generate a `cpu.prof` file. Inspect that file with `go tool pprof cpu.prof`.

`make test-san` will run two very slow tests: one with race-condition detection, the second with the memory sanitizer. This requires Clang to be installed with a complete `compiler-rt` package; Alpine distro packages won't work. `-msan` is only supported on Linux amd64/arm64; see `go help build` for more info.

