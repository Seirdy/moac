#!/bin/sh
# install all linters and formatters to GOBIN.

set -e -u -x

# golangci-lint is installed from master
linters='github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0
github.com/praetorian-inc/gokart@latest
github.com/mrtazz/checkmake/cmd/checkmake@latest
github.com/quasilyte/go-consistent@latest
mvdan.cc/sh/v3/cmd/shfmt@latest
github.com/sonatype-nexus-community/nancy@latest
github.com/fe3dback/go-arch-lint@latest'

formatters='mvdan.cc/gofumpt@latest
golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
github.com/mdigger/goldmark-formatter/cmd/mdfmt@latest'

installme="$linters
$formatters"

set +u

# if we're running in CI, only get linters. builds.sr.ht exports JOB_ID
[ -n "$JOB_ID" ] && installme="$linters"

set -u

tmp_dir="$(mktemp -d)"
workdir="$PWD"
cd "$tmp_dir"
# first download the linters in parallel, then compile them one by one
# since compilation uses parallelism
echo "$installme" | xargs -n1 -P2 go get -d
cd "$workdir"

go_install() {
	CGO_ENABLED=0 go install -trimpath -mod=readonly -buildmode=exe -ldflags '-w -s -linkmode=internal' "$1"
}

for tool in $installme; do
	go_install "$tool"
done

echo 'installed linters successfully'
