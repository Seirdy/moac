.POSIX:

MOAC_BIN = moac
MOAC_PWGEN_BIN = moac-pwgen
BINS= $(MOAC_BIN) $(MOAC_PWGEN_BIN)

SHARED_SRC = Makefile *.go entropy/*.go internal/*/*.go
MOAC_SRC = cmd/moac/*.go
MOAC_PWGEN_SRC = pwgen/*.go cmd/moac-pwgen/*.go
SRC = $(SHARED_SRC) $(MOAC_EXCLUSIVE_SRC) $(MOAC_PWGEN_EXCLUSIVE_SRC)

CGO_ENABLED ?= 0
GOPATH ?= `$(GO) env GOPATH`
GOBIN ?= $(GOPATH)/bin
COVERPKG = .,./entropy,./pwgen,./internal/slicing


GO ?= go
GOLANGCI_LINT ?= $(GOBIN)/golangci-lint
GOKART ?= $(GOBIN)/gokart
GOKART_FLAGS ?= -g

CMD = build
ARGS =
TAG = `git describe --abbrev=0 --tags`
REVISION = `git rev-parse --short HEAD`
VERSION = $(TAG)-$(REVISION)

PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin
DATAROOTDIR ?= $(PREFIX)/share
MANDIR ?= $(DATAROOTDIR)/man
ZSHCOMPDIR ?= $(DATAROOTDIR)/zsh/site-functions

# general build flags
LINKMODE = internal
# extldflags is ignored unless you use one of the cgo options at the bottom
GO_LDFLAGS += -w -s -X git.sr.ht/~seirdy/moac/internal/version.version='$(VERSION)' -linkmode='$(LINKMODE)' -extldflags \"$(LDFLAGS)\"
BUILDMODE ?= default
GO_BUILDFLAGS += -trimpath -mod=readonly -gcflags=-trimpath=$(GOPATH) -asmflags=-trimpath=$(GOPATH) -buildmode=$(BUILDMODE) -ldflags='$(GO_LDFLAGS)'
TESTFLAGS ?= # -msan, -race, coverage, etc.

all: build doc

golangci-lint: $(SRC)
	$(GOLANGCI_LINT) run
gokart-lint: $(SRC)
	$(GOKART) scan $(GOKART_FLAGS) ./...

lint: golangci-lint gokart-lint

.base: $(SRC)
	CC=$(CC) CCLD=$(CCLD) CGO_CFLAGS="$(CFLAGS)" CGO_ENABLED=$(CGO_ENABLED) $(GO) $(CMD) $(GO_BUILDFLAGS) $(ARGS)

$(MOAC_BIN): $(SHARED_SRC) $(MOAC_SRC)
	@$(MAKE) GO_BUILDFLAGS="$(GO_BUILDFLAGS) -o $(MOAC_BIN)" CMD=build ARGS=./cmd/moac .base
$(MOAC_PWGEN_BIN): $(SHARED_SRC) $(MOAC_PWGEN_SRC)
	@$(MAKE) GO_BUILDFLAGS="$(GO_BUILDFLAGS) -o $(MOAC_PWGEN_BIN)" CMD=build ARGS=./cmd/moac-pwgen .base

build: $(BINS)

clean:
	$(GO) clean -testcache
	rm -f $(BINS) doc/*.1 ./coverage.out

test:
	@$(MAKE) CMD="test" GO_BUILDFLAGS="$(GO_BUILDFLAGS)" ARGS="$(TESTFLAGS) ./..." .base

test-cov:
	@$(MAKE) TESTFLAGS="-coverpkg=$(COVERPKG) -coverprofile=coverage.out" test
	$(GO) tool cover -func=coverage.out

fmt:
	fieldalignment -fix ./...
	gofumpt -s -w .

pre-commit: fmt lint test

doc/moac.1: doc/moac.1.scd
	scdoc < doc/moac.1.scd > doc/moac.1
doc/moac-pwgen.1: doc/moac-pwgen.1.scd
	scdoc < doc/moac-pwgen.1.scd > doc/moac-pwgen.1

doc: doc/moac.1 doc/moac-pwgen.1

install: all
	mkdir -p \
		$(DESTDIR)$(BINDIR) \
		$(DESTDIR)$(MANDIR)/man1 \
		$(DESTDIR)$(ZSHCOMPDIR)
	cp -f $(BINS) $(DESTDIR)$(BINDIR)
	chmod 755 $(DESTDIR)$(BINDIR)/$(MOAC_BIN)
	chmod 755 $(DESTDIR)$(BINDIR)/$(MOAC_PWGEN_BIN)
	cp -f doc/*.1 $(DESTDIR)$(MANDIR)/man1
	chmod 644 $(DESTDIR)$(MANDIR)/man1/moac.1
	chmod 644 $(DESTDIR)$(MANDIR)/man1/moac-pwgen.1
	cp -f completions/zsh/_* $(DESTDIR)$(ZSHCOMPDIR)
	chmod 644 $(DESTDIR)$(ZSHCOMPDIR)/_moac $(DESTDIR)$(ZSHCOMPDIR)/_moac-pwgen

# =================================================================================

# everything below this line requires CGO + Clang. Building with CGO allows a few
# extra goodies:
# 	- static-pie binaries (note that go puts the heap at a fixed address
# 	  anyway and this isn't useful without CGO)
# 	- msan and race detection (msan requires clang)
# 	- support for platforms that require CGO like OpenBSD

# moac doesn't really need CGO outside platforms like FreeBSD but this Makefile is
# just a template that I use for all my Go projects. Besides, it should be safe to build it with CGO

# if building with CGO, turn on some hardening
CC = clang
CCLD = lld
CFLAGS = -O2 -fno-semantic-interposition -g -pipe -Wp,-D_FORTIFY_SOURCE=2 -Wp,-D_GLIBCXX_ASSERTIONS -fexceptions -fstack-protector-all -m64 -fasynchronous-unwind-tables -fstack-clash-protection -fcf-protection=full -ffunction-sections -fdata-sections -ftrivial-auto-var-init=zero -enable-trivial-auto-var-init-zero-knowing-it-will-be-removed-from-clang
LDFLAGS = -Wl,-z,relro,-z,now,-z,noexecstack -Wl,--as-needed -Wl,-E -Wl,--gc-sections
# on openbsd, set this to "exe" or nothing
BUILDMODE_CGO = pie

# on ARMv8, you can switch safe-stack to shadow-call-stack
# on Alpine, set this to cfi since compiler-rt isn't built properly.
# openbsd doesn't support either
EXTRA_SANITIZERS ?= cfi
CFI = -flto=thin -fsanitize=$(EXTRA_SANITIZERS)
CFLAGS_CFI = $(CFLAGS) $(CFI) -fvisibility=hidden -fpic -fpie
LDFLAGS_CFI = $(LDFLAGS) $(CFI) -pie
EXTRA_LDFLAGS =

# Test with thread and memory sanitizers; needs associated libclang_rt libs.
.test-cgo:
	@$(MAKE) CGO_ENABLED=1 CFLAGS="$(CFLAGS_CFI)" LINKMODE=external test

test-race:
	@$(MAKE) TESTFLAGS='-race' .test-cgo

# test-msan does not work on alpine (its compiler-rt lacks msan)
# but it works on fedora and void-musl.
test-msan:
	@$(MAKE) TESTFLAGS='-msan' BUILDMODE=$(BUILDMODE_CGO) .test-cgo

test-san: test-race test-msan

.build-cgo-base:
	@$(MAKE) CFLAGS="$(CFLAGS_CFI)" CGO_ENABLED=1 LINKMODE=external BUILDMODE=$(BUILDMODE_CGO) LDFLAGS="$(LDFLAGS_CFI) $(EXTRA_LDFLAGS)" build

build-cgo:
	@$(MAKE) .build-cgo-base

# build a static-pie binary with sanitizers for CFI and either
# safe-stack (x86_64) or shadow-call-stack (ARMv8)
# the below should be run on a musl-based toolchain; works on Alpine or Void-musl
# Tends to cause crashes when linking with glibc
# alpine users should disable safe-stack since the alpine compiler-rt package is incomplete
build-cgo-static:
	@$(MAKE) EXTRA_LDFLAGS=-static-pie build-cgo

.PHONY: all clean doc lint fmt pre-commit test test-race test-msan test-san test-cov build build-cgo build-cgo-static install
