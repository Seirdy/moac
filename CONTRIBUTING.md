Contributing to MOAC
====================

I try to accommodate everyone's workflow. Here's ways to contribute and how, in my order of preference:

Security vulnerabilities, non-public contact
--------------------------------------------

If you want to connect with me directly: my email address is at the bottom of the man pages or in the commit logs. Vulnerability disclosures should be PGP-encrypted using the PGP key [1E892DB2A5F84479](https://seirdy.one/publickey.asc).

Alternatively, send an encrypted message on Matrix to `@seirdy:envs.net`

Bug reports and TODOs
---------------------

Preferred and canonical ticket tracker: <https://todo.sr.ht/~seirdy/MOAC>. Send an email to [~seirdy/moac@todo.sr.ht](mailto:~seirdy/MOAC@todo.sr.ht) to automatically file a bug, no account needed. The tracker might also have some tickets labeled "good first issue" ideal for contributors with less experience.

I also check issues in the GitHub, GitLab, and Codeberg mirrors linked at the top of the README, if you prefer. No matter which option you choose, your bug gets emailed to me.

Patches, questions, and feature requests
----------------------------------------

Preferred location: <https://lists.sr.ht/~seirdy/moac>. Send emails and patches to [~seirdy/moac@lists.sr.ht](mailto:~seirdy/moac@lists.sr.ht). I also check the GitHub, GitLab, and Codeberg mirrors for issues and PRs.

### Coding standards

Contributions don't need to follow these standards to be useful. If a useful patch doesn't pass the below checks, I might clean it up myself.

This project uses `gofumpt` and `fieldalignment` for formatting; it uses `golangci-lint`, `gokart`, and `checkmake` for linting. You can install all of them to your `GOBIN` like so:

```sh
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/praetorian-inc/gokart@latest
go install github.com/mrtazz/checkmake/cmd/checkmake@latest
go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
go install github.com/mvdan/gofumpt@latest
```

Run `make fmt` to format code, `make lint` to run the linters, and `make test` to run unit tests. `make pre-commit` runs all three.

The linters are very opinionated. If you find this annoying, you can send your patch anyway; I'll clean it up if it looks useful.

Everything possible should have tests. If a line returning an error is impossible to reach and thus uncovered, replace it with a panic. Any uncovered line that isn't a panic is in need of a test.

