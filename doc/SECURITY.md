MOAC Security
=============

Users place a very high degree of trust in password generators and evaluators. MOAC must therefore meet a high bar for security standards.

Security requirements
---------------------

Security features users can and cannot expect:

- Password generation, the only component that features non-determinism, exclusively uses the CSPRNG offered by the Go stdlib's `crypto/rand` package. Check the GoDoc for `crypto/rand` to see its security standards.
- Entropy measurement is based solely on password length and charsets used; it does not take into account any other characteristics such as dictionary words, repetition, etc. Entropy measurement was designed under the assumption that the measured passwords were randomly generated.
- Password strength metrics depend only on physical laws, never needing to be updated to account for advancements in computing power.
- Password-crackability metrics do not assume the presence of a key-derivation function or key stretching/strengthening. Making fewer assumptions helps maintain simplicity and applicability to the widest range of threat models.
- Simplicity: MOAC should have a limited scope (password analysis and generation) and size (<1k Go SLOC, excluding tests). This isn't technically a security requirement, but it does keep attack surface low and reduce room for bugs.

Third-party dependencies
------------------------

MOAC is split into a library and two CLI utilities. The library has no third-party dependencies; the CLI utilities have a few. A CI job scans these dependencies and indirect dependencies against Sonatype's OSS Index on every push.

Builds
------

- MOAC should never require any use of C libraries or dynamic linking; binaries can be 100% Go-based (with the exception of OpenBSD binaries, which use CGO for syscalls as of Go 1.16) to ensure a high level of memory safety. Depending on build flags, it's still possible to use CGO (e.g. if using `-buildmode=pie`), but it should never be a requirement.
- MOAC supports reproducible builds that contain bit-for-bit identical binaries for a given Go toolchain.

Checks and enforcement
----------------------

Every push triggers CI jobs that run several tests in VMs.

- Every reachable, non-deprecated statement in the library should be covered by tests. Mutation testing should reveal a mutation score above 0.7.
- Password generation is tested especially heavily, with thousands of pairwise test-cases assembled from combinations of valid parameters covering known edge cases.
- Furthermore, each pairwise case undergoes gorilla testing to account for the non-determinism of password generation. In CI jobs, this includes 512 repetitions in an OpenBSD VM and 128 repetitions in an Alpine Linux VM.
- A Fedora Linux VM runs tests with Go's memory sanitizer and race detector enabled.
- Every push undergoes strict static analysis that includes GoKart and every relevant linter in golangci-lint. Check those projects to see which vulnerabilities they cover.

Support policy
--------------

MOAC follows a rolling release system, and does not backport fixes to previous versions. Deprecated code may receive security fixes, but won't be held to the same testing standards.

