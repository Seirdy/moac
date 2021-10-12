MOAC Security
=============

Security requirements
---------------------

- Password generation is the only source of non-determinism from random-number generation. It exclusively uses CSPRNG offered by the `crypto/rand` package from the Go standard library.
- Entropy measurement is based solely on password length and charsets used; it does not take into account any other characteristics such as dictionary words, repetition, etc. Entropy measurement was designed under the assumption that the measured passwords were randomly generated.
- Password strength metrics depends only on physical laws, never needing to be updated to account for advancements in computing power.
- Password-crackability metrics do not assume the presence of a key-derivation function or key stretching/strengthening. Making fewer assumptions helps maintain simplicity and applicability to the widest range of threat models.
- Simplicity: MOAC should have a limited scope (password analysis and generation) and size (<1k SLOC, excluding tests). This isn't technically a security requirement, but it does keep attack surface low and reduce room for bugs.

Dependencies
------------

- The MOAC library has no third-party dependencies. The CLI utilities' third-party dependencies are limited to official libraries from `golang.org/x/`, a simple `getopts`-like flag parser, and a testing library. Some of these include indirect dependencies that are only used for testing the dependencies and are not included in the final binaries.
- A CI job scans dependencies against Sonatype's OSS Index on every push.

Builds
------

- MOAC should never require any use of C libraries or dynamic linking; binaries can be 100% Go-based (with the exception of OpenBSD, which uses CGO for all Go binaries as of Go 1.16) to ensure a high level of memory safety. Depending on build flags, it's still possible to use CGO (e.g. if using `-buildmode=pie`), but it should never be a requirement.
- MOAC supports reproducible builds that contain bit-for-bit identical binaries for a given Go toolchain.

Checks and enforcement
----------------------

Every push triggers CI jobs that run several tests in VMs

- Every reachable, non-deprecated statement in the library should be covered by tests. Mutation testing should reveal a mutation score above 0.7.
- Password generation is tested especially heavily, with thousands of test-cases assembled from combinations of valid parameters covering known edge cases. Furthermore, each test case is tested 512 times in an OpenBSD VM and 128 times in an Alpine VM due to the non-determinism of password generation.
- One VM runs tests with Go's memory sanitizer and race detector.
- Every push undergoes strict static analysis that includes GoKart and every relevant linter in golangci-lint. Check those projects to see which vulnerabilities they cover.

