#!/bin/sh

linters='github.com/golangci/golangci-lint/cmd/golangci-lint@latest github.com/praetorian-inc/gokart@latest github.com/mrtazz/checkmake/cmd/checkmake@latest'

# first download the linters in parallel, then compile them one by one
echo "$linters" | tr ' ' '\n' | xargs -P2 -L1 go get -d

for linter in $linters; do
	go install "$linter"
done
