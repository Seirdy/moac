#!/bin/sh
# Install Go, set linker/CC to lld/clang
# Works for glibc-based Linux distributions.

set -e -u

# Binaries made with `make dist-reprod` are reproducible for a given Go installation.
# All builds with the below Go installation, for example, should have the same checksums.

go_version=1.17.3
go_hash=550f9845451c0c94be679faf116291e7807a8d78b43149f9506c1b15eb89008c
go_tarball="go$go_version.linux-amd64.tar.gz"

curl -sSLo "$go_tarball" "https://dl.google.com/go/go$go_version.linux-amd64.tar.gz"
found_hash="$(sha256sum "$go_tarball" | cut -d' ' -f1)"

if [ "$found_hash" != "$go_hash" ]; then
	echo "Checksum mismatch: $found_hash"
	exit 1
fi

sudo tar -C /usr/local -xzf "$go_tarball"

# Below is just for the sanitizers, not for setting up a reproducible build
sudo ln -sf /usr/bin/ld.lld /usr/bin/ld
sudo ln -sf /usr/bin/clang /usr/bin/cc
