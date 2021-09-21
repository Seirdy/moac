#!/bin/sh
# install Go, set linker/CC to lld/clang

set -e -u

go_version=1.17.1
go_hash=dab7d9c34361dc21ec237d584590d72500652e7c909bf082758fb63064fca0ef

go_tarball="go$go_version.linux-amd64.tar.gz"
curl -sSLo "$go_tarball" "https://storage.googleapis.com/golang/go$go_version.linux-amd64.tar.gz"
found_hash="$(sha256sum "$go_tarball" | cut -d' ' -f1)"

if [ "$found_hash" != "$go_hash" ]; then
	echo "Checksum mismatch: $found_hash"
	exit 1
fi

sudo tar -C /usr/local -xzf "$go_tarball"
sudo ln -sf /usr/bin/ld.lld /usr/bin/ld
sudo ln -sf /usr/bin/clang /usr/bin/cc
