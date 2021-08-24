#!/bin/sh
set -e

sudo alternatives --set ld /usr/bin/ld.lld
ln -sf /usr/bin/clang-1* /usr/bin/cc
