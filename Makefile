.POSIX:

BIN = moac
CGO_ENABLED ?= 0
GOPATH ?= $(shell $(GO) env GOPATH)
GOBIN ?= $(GOPATH)/bin
SRC = *.go entropy/*.go cmd/moac/*.go

GO ?= go
GOLANGCI_LINT ?= $(GOBIN)/golangci-lint
GOKART ?= $(GOBIN)/gokart
GOKART_FLAGS ?= -g

# general build flags
GO_LDFLAGS += "-w -s"
GO_BUILDFLAGS += -trimpath -mod=readonly -gcflags="-trimpath=$(GOPATH)" -asmflags="-trimpath=$(GOPATH)"

default:
	$(MAKE) clean build

golangci-lint: $(SRC)
	$(GOLANGCI_LINT) run
gokart-lint: $(SRC)
	$(GOKART) scan $(GOKART_FLAGS) .
	$(GOKART) scan $(GOKART_FLAGS) ./entropy
	$(GOKART) scan $(GOKART_FLAGS) ./cmd/moac

lint: golangci-lint gokart-lint

$(BIN): $(SRC)
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build $(GO_BUILDFLAGS) -buildmode=exe -ldflags $(GO_LDFLAGS) -o $(BIN) ./cmd/moac

build: $(BIN)

clean:
	$(GO) clean
	rm -f ./$(BIN) ./coverage.out

test: $(SRC)
	CGO_ENABLED=0 $(GO) test $(GO_BUILDFLAGS) ./...

test-cov: $(SRC)
	CGO_ENABLED=0 $(GO) test $(GO_BUILDFLAGS) -coverpkg=.,./entropy -coverprofile=coverage.out ./...

# =================================================================================

# everything below this line requires CGO + Clang. Building with CGO allows a few
# extra goodies:
# 	- PIE and static-pie binaries
# 	- msan and race detection

# if building with CGO, turn on some hardening
CC = clang
CCLD = lld
CFLAGS += -O2 -fno-semantic-interposition -g -pipe -Wp,-D_FORTIFY_SOURCE=2 -Wp,-D_GLIBCXX_ASSERTIONS -fexceptions -fstack-protector-all -m64 -fasynchronous-unwind-tables -fstack-clash-protection -fcf-protection=full -ffunction-sections -fdata-sections -ftrivial-auto-var-init=zero -enable-trivial-auto-var-init-zero-knowing-it-will-be-removed-from-clang
LDFLAGS += -Wl,-z,relro,-z,now,-z,noexecstack -Wl,--as-needed -Wl,-E -Wl,--gc-sections
GO_LDFLAGS_CGO += "-w -s -linkmode=external -extldflags '$(LDFLAGS)'"
# on openbsd, set this to "exe" or nothing
BUILDMODE_CGO ?= pie

# for testing with clang+msan+CFI+safe-stack/shadow-stack and release builds
CFLAGS_LTO_PIE += $(CFLAGS) -flto=thin -fvisibility=hidden -fpic -fpie
EXTRA_SANITIZERS ?= cfi
CFLAGS_CFI += $(CFLAGS_LTO_PIE) -fsanitize=$(EXTRA_SANITIZERS)
LDFLAGS_CFI += $(LDFLAGS) -flto=thin -fsanitize=$(EXTRA_SANITIZERS) -pie
GO_LDFLAGS_CFI += "-w -s -linkmode=external -extldflags '$(LDFLAGS_CFI)'"

# for release builds, with Clang+CFI sanitization, static-pie linked
# on ARMv8, you can switch safe-stack to shadow-call-stack
# on Alpine, set this to cfi since compiler-rt isn't built properly.
RELEASE_SANITIZERS ?= cfi,safe-stack
LDFLAGS_RELEASE += $(LDFLAGS) -flto=thin -fsanitize=$(RELEASE_SANITIZERS) -static-pie
CFLAGS_RELEASE += $(CFLAGS_LTO_PIE) -fsanitize=$(RELEASE_SANITIZERS)
GO_LDFLAGS_RELEASE += "-w -s -linkmode=external -extldflags '$(LDFLAGS_RELEASE)'"

# Test with thread and memory sanitizers; needs associated libclang_rt libs.
test-race: $(SRC)
	CC=$(CC) CCLD=$(CCLD) CGO_CFLAGS="$(CFLAGS)" $(GO) test $(GO_BUILDFLAGS) -race -ldflags=$(GO_LDFLAGS_CGO) -coverpkg=.,./entropy .

# test-msan does not work on alpine (its compiler-rt lacks msan)
# but it works on fedora and void-musl.
test-msan: $(SRC)
	CC=clang CCLD=lld CGO_CFLAGS="$(CFLAGS_CFI)" $(GO) test $(GO_BUILDFLAGS) -buildmode=$(BUILDMODE_CGO) -msan -ldflags=$(GO_LDFLAGS_CFI) .

test-san: test-race test-msan

build-cgo: $(SRC)
	CC=$(CC) CCLD=$(CCLD) CGO_CFLAGS="$(CFLAGS_CFI)" $(GO) build $(GO_BUILDFLAGS) -buildmode=$(BUILDMODE_CGO) -ldflags=$(GO_LDFLAGS_CFI) -o $(BIN) ./cmd/moac

# build a static-pie binary with sanitizers for CFI and either
# safe-stack (x86_64) or shadow-call-stack (ARMv8)
# the below should be run on a musl-based toolchain; works on Alpine or Void-musl
# Tends to cause crashes when linking with glibc
build-cgo-static: $(SRC)
	CC=$(CC) CCLD=$(CCLD) CGO_CFLAGS="$(CFLAGS_RELEASE)" $(GO) build $(GO_BUILDFLAGS) -buildmode=$(BUILDMODE_CGO) -ldflags=$(GO_LDFLAGS_RELEASE) -o $(BIN) ./cmd/moac

.PHONY: all lint test test-race test-msan test-san test-cov build build-cgo build-cgo-static
