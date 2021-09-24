.POSIX:

# built binary names
MOAC_BIN = moac
MOAC_PWGEN_BIN = moac-pwgen
BINS= $(MOAC_BIN) $(MOAC_PWGEN_BIN)

# source files
SHARED_SRC = Makefile *.go entropy/*.go internal/*/*.go
MOAC_SRC = cmd/moac/*.go
MOAC_PWGEN_SRC = pwgen/*.go cmd/moac-pwgen/*.go
SRC = $(SHARED_SRC) $(MOAC_EXCLUSIVE_SRC) $(MOAC_PWGEN_EXCLUSIVE_SRC)
COVERPKG = .,./entropy,./pwgen,./charsets,./internal/bounds

# go's own envvars
CGO_ENABLED ?= 0
GOPATH ?= `$(GO) env GOPATH`
GOBIN ?= $(GOPATH)/bin
GOOS ?= `$(GO) env GOOS`
GOARCH ?= `$(GO) env GOARCH`
CGO_CFLAGS +=  $(CFLAGS)

# paths to executables this Makefile will use
GO ?= go
GOLANGCI_LINT ?= $(GOBIN)/golangci-lint
GOKART ?= $(GOBIN)/gokart
CHECKMAKE ?= $(GOBIN)/checkmake
GOFUMPT ?= $(GOBIN)/gofumpt
FIELDALIGNMENT ?= $(GOBIN)/fieldalignment

# change this on freebsd/openbsd
SHA256 ?= sha256sum

# version identifier to embed in binaries
TAG = `git describe --abbrev=0 --tags`
REVISION = `git rev-parse --short HEAD`
VERSION = $(TAG)-$(REVISION)

# install destinations
PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin
DATAROOTDIR ?= $(PREFIX)/share
MANDIR ?= $(DATAROOTDIR)/man
ZSHCOMPDIR ?= $(DATAROOTDIR)/zsh/site-functions

# general build flags
LINKMODE = internal
# extldflags is ignored unless you use one of the cgo options at the bottom
DEFAULT_GO_LDFLAGS = -w -X git.sr.ht/~seirdy/moac/internal/version.version="$(VERSION)" -linkmode=$(LINKMODE) -extldflags \"$(LDFLAGS)\"
GO_LDFLAGS += $(DEFAULT_GO_LDFLAGS)
BUILDMODE ?= default
GO_BUILDFLAGS += -trimpath -mod=readonly -buildmode=$(BUILDMODE) -ldflags '$(GO_LDFLAGS)'
TESTFLAGS ?= # -msan, -race, coverage, etc.

# used internally
CMD = build
ARGS =

all: build doc

golangci-lint: $(SRC)
	$(GOLANGCI_LINT) run
gokart-lint: $(SRC)
	$(GOKART) scan -g ./...
checkmake: Makefile
	$(CHECKMAKE) Makefile

lint: golangci-lint gokart-lint checkmake

# every task in this makefile except "clean" just calls .base with different vars
# instead of invoking "$(GO)" directly
.base: $(SRC)
	$(GO) $(CMD) $(GO_BUILDFLAGS) $(ARGS)

$(MOAC_BIN): $(SHARED_SRC) $(MOAC_SRC)
	$(MAKE) GO_BUILDFLAGS="$(GO_BUILDFLAGS) -o $(MOAC_BIN)" CMD=build ARGS=./cmd/moac .base
$(MOAC_PWGEN_BIN): $(SHARED_SRC) $(MOAC_PWGEN_SRC)
	$(MAKE) GO_BUILDFLAGS="$(GO_BUILDFLAGS) -o $(MOAC_PWGEN_BIN)" CMD=build ARGS=./cmd/moac-pwgen .base

build: $(BINS)

.clean-bins:
	rm -f $(BINS)

clean: .clean-bins clean-san-bins
	$(GO) clean -testcache
	rm -rf doc/*.1 ./coverage.out $(DIST)

test:
	@$(MAKE) CMD="test" GO_BUILDFLAGS="$(GO_BUILDFLAGS)" ARGS="$(TESTFLAGS) ./..." .base

test-cov:
	@$(MAKE) TESTFLAGS="-coverpkg=$(COVERPKG) -coverprofile=coverage.out" test
	$(GO) tool cover -func=coverage.out

fmt:
	$(FIELDALIGNMENT) -fix ./...
	$(GOFUMPT) -s -w .

pre-commit: fmt lint test

doc/moac.1: doc/moac.1.scd
	scdoc < doc/moac.1.scd > doc/moac.1
doc/moac-pwgen.1: doc/moac-pwgen.1.scd
	scdoc < doc/moac-pwgen.1.scd > doc/moac-pwgen.1

doc: doc/moac.1 doc/moac-pwgen.1

# final install jobs include these two targets
INSTALL_SHARE = install-man install-completion

install-bin: build
	mkdir -p $(DESTDIR)$(BINDIR)
	cp -f $(BINS) $(DESTDIR)$(BINDIR)
	chmod 755 $(DESTDIR)$(BINDIR)/$(MOAC_BIN) $(DESTDIR)$(BINDIR)/$(MOAC_PWGEN_BIN)
install-bin-strip:
	$(MAKE) GO_LDFLAGS='$(GO_LDFLAGS) -s' install-bin
install-man: doc
	mkdir -p  $(DESTDIR)$(MANDIR)/man1
	cp -f doc/*.1 $(DESTDIR)$(MANDIR)/man1
	chmod 644 $(DESTDIR)$(MANDIR)/man1/moac.1 $(DESTDIR)$(MANDIR)/man1/moac-pwgen.1
install-completion:
	mkdir -p  $(DESTDIR)$(ZSHCOMPDIR)
	cp -f completions/zsh/_* $(DESTDIR)$(ZSHCOMPDIR)
	chmod 644 $(DESTDIR)$(ZSHCOMPDIR)/_moac $(DESTDIR)$(ZSHCOMPDIR)/_moac-pwgen

install: install-bin $(INSTALL_SHARE)
install-strip: install-bin-strip $(INSTALL_SHARE)

uninstall:
	rm -f \
		$(DESTDIR)$(BINDIR)/$(MOAC_BIN) $(DESTDIR)$(BINDIR)/$(MOAC_PWGEN_BIN) \
		$(DESTDIR)$(MANDIR)/man1/moac.1 $(DESTDIR)$(MANDIR)/man1/moac-pwgen.1 \
		$(DESTDIR)$(ZSHCOMPDIR)/_moac $(DESTDIR)$(ZSHCOMPDIR)/_moac-pwgen

# =================================================================================
# Build tarballs containing reproducible builds

PLATFORM_ID = $(GOOS)-$(GOARCH)
RELNAME = moac-$(VERSION)-$(PLATFORM_ID)
# allow excluding the git version from the archive name
# this lets the archive name be deterministic, which is useful in CI
# because sourcehut artifact names are interpreted literally.
ARCHIVE_PREFIX ?= moac-$(VERSION) # override ARCHIVE_PREFIX in CI
ARCHIVE_NAME ?= $(ARCHIVE_PREFIX)-$(PLATFORM_ID)
DIST ?= dist
DIST_LOCATION=$(DIST)/$(RELNAME)

dist:
	DESTDIR=$(DIST)/$(RELNAME) $(MAKE) install-strip
	@$(SHA256) $(DIST)/$(RELNAME)/$(BINDIR)/*
	@tar czf "$(DIST)/$(ARCHIVE_NAME).tar.gz" -C $(DIST)/ $(RELNAME)
	@rm -rf $(DIST)/$(RELNAME)

# For reproducible builds, throw out the non-deterministic Build ID
# Install Go to /usr/local for builds with that Go toolchain to be reproducible.
dist-reprod:
	$(MAKE) GO_LDFLAGS='-buildid= $(DEFAULT_GO_LDFLAGS)' .clean-bins dist

dist-multiarch:
	@GOARCH=amd64 $(MAKE) dist-reprod
	@GOARCH=arm64 $(MAKE) dist-reprod
	@GOARCH=arm $(MAKE) dist-reprod
	@GOARCH=386 $(MAKE) dist-reprod

# This builds for Linux and FreeBSD, but not OpenBSD; OpenBSD bins should
# be built with CGO which makes reproducible cross-compilation a bit messy
dist-linux-freebsd:
	@GOOS=linux $(MAKE) dist-multiarch
	@GOOS=freebsd $(MAKE) dist-multiarch


# =================================================================================

# everything below this line requires CGO + Clang. Building with CGO allows a few
# extra goodies:
# 	- static-pie binaries (note that go puts the heap at a fixed address
# 	  anyway and this isn't useful without CGO, but some platforms enforce PIE-ness)
# 	- msan and race detection (msan requires clang)
# 	- support for platforms that require CGO like OpenBSD

# moac doesn't really need CGO outside platforms like OpenBSD but this Makefile is
# just a template that I use for all my Go projects.

# if building with CGO, turn on some hardening
CC = clang
CCLD = lld
CFLAGS += -O2 -fno-semantic-interposition -g -pipe -Wp,-D_FORTIFY_SOURCE=2 -Wp,-D_GLIBCXX_ASSERTIONS -fexceptions -fstack-protector-all -m64 -fasynchronous-unwind-tables -fstack-clash-protection -fcf-protection=full -ffunction-sections -fdata-sections -ftrivial-auto-var-init=zero -enable-trivial-auto-var-init-zero-knowing-it-will-be-removed-from-clang
LDFLAGS += -Wl,-z,relro,-z,now,-z,noexecstack,--as-needed,-E,--gc-sections
# on openbsd, set this to "exe" or nothing
BUILDMODE_CGO = pie

# on ARMv8, you can switch safe-stack to shadow-call-stack
# on Alpine, set this to cfi since compiler-rt isn't built properly.
# openbsd doesn't support either
EXTRA_SANITIZERS ?= cfi
CFI = -flto=thin -fsanitize=$(EXTRA_SANITIZERS)
CFLAGS_CFI = $(CFLAGS) $(CFI) -fvisibility=hidden -fpic -fpie
LDFLAGS_CFI = $(LDFLAGS) $(CFI) -pie

# shared across regular, msan, and race CGO builds/tests
.build-cgo-base:
	@CC="$(CC)" CCLD="$(CCLD)" CFLAGS="$(CFLAGS_CFI)" LDFLAGS="$(LDFLAGS)" CGO_CFLAGS="$(CFLAGS_CFI)" $(MAKE) CGO_ENABLED=1 LINKMODE=external $(CMD)

build-cgo:
	@$(MAKE) BUILDMODE=$(BUILDMODE_CGO) CFLAGS="$(CFLAGS_CFI)" LDFLAGS="$(LDFLAGS_CFI)" .build-cgo-base

build-msan:
	@GO_BUILDFLAGS='-msan' $(MAKE) .build-cgo-base

build-race:
	@GO_BUILDFLAGS='-race' $(MAKE) CGO_ENABLED=1 .build-cgo-base

build-san:
	@$(MAKE) MOAC_BIN=$(MOAC_BIN)-msan MOAC_PWGEN_BIN=$(MOAC_PWGEN_BIN)-msan build-msan
	@$(MAKE) MOAC_BIN=$(MOAC_BIN)-race MOAC_PWGEN_BIN=$(MOAC_PWGEN_BIN)-race build-race

# cleans just the artifacts produced by build-san
clean-san-bins:
	@$(MAKE) CMD=.clean-bins build-san

test-race:
	@$(MAKE) CMD='test' build-race

help:
	grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST)

# test-msan does not work on alpine (its compiler-rt lacks msan)
# but it works on fedora and void-musl.
test-msan:
	@$(MAKE) CMD='test' build-msan

test-san: test-race test-msan

# build a static-pie binary with sanitizers for CFI and either
# safe-stack (x86_64) or shadow-call-stack (ARMv8)
# the below should be run on a musl-based toolchain; works on Alpine or Void-musl
# Tends to cause crashes when linking with glibc
# alpine users should disable safe-stack since the alpine compiler-rt package is incomplete
build-cgo-static:
	@$(MAKE) LDFLAGS_CFI='$(LDFLAGS) $(CFI) -static-pie' build-cgo

.PHONY: test test-race test-msan test-san test-cov
.PHONY: build .build-cgo-base build-cgo build-cgo-static build-msan build-race
.PHONY: install-bin install-man install-completion install-bin-strip install-strip install
.PHONY: all clean .clean-bins doc lint fmt pre-commit pre-push uninstall
.PHONY: dist dist-reprod dist-multiarch dist-linux-freebsd
