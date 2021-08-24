.POSIX:

BIN = moac
CGO_ENABLED ?= 0
GOPATH ?= $(shell $(GO) env GOPATH)
GOBIN ?= $(GOPATH)/bin

GO ?= go
GOLANGCI_LINT ?= $(GOBIN)/golangci-lint
GOKART ?= $(GOBIN)/gokart
GOKART_FLAGS ?= -g

# general build flags
GO_BUILDFLAGS += -trimpath -mod=readonly -gcflags="-trimpath=$(GOPATH)" -asmflags="-trimpath=$(GOPATH)"

default:
	$(MAKE) clean build

golangci-lint:
	$(GOLANGCI_LINT) run
gokart-lint:
	$(GOKART) scan $(GOKART_FLAGS) .
	$(GOKART) scan $(GOKART_FLAGS) ./entropy
	$(GOKART) scan $(GOKART_FLAGS) ./cmd/moac

lint: golangci-lint gokart-lint

$(BIN):
	CC=$(CC) CCLD=$(CCLD) CGO_ENABLED=$(CGO_ENABLED) $(GO) build $(GO_BUILDFLAGS) -o $(BIN) ./cmd/moac/

build: $(BIN)

clean:
	$(GO) clean
	rm -f ./moac ./coverage.out

test:
	CGO_ENABLED=0 $(GO) test $(GO_BUILDFLAGS) ./...

test-cov:
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
GO_LDFLAGS += "-w -s -linkmode=external -extldflags '$(LDFLAGS)'"

# for testing with clang+msan+CFI+safe-stack/shadow-stack and release builds
CFLAGS_LTO_PIE += $(CFLAGS) -flto=thin -fvisibility=hidden -fpic -fpie
EXTRA_SANITIZERS ?= cfi
CFLAGS_CFI += $(CFLAGS_LTO_PIE) -fsanitize=$(EXTRA_SANITIZERS)
LDFLAGS_CFI += $(LDFLAGS) -flto=thin -fsanitize=cfi -pie
GO_LDFLAGS_CFI += "-w -s -linkmode=external -extldflags '$(LDFLAGS_CFI)'"

# for release builds, with Clang+CFI sanitization, static-pie linked
# on ARMv8, you can switch safe-stack to shadow-call-stack
# on Alpine, set this to cfi since compiler-rt isn't built properly.
RELEASE_SANITIZERS ?= cfi,safe-stack
LDFLAGS_RELEASE += $(LDFLAGS) -flto=thin -fsanitize=$(RELEASE_SANITIZERS) -static-pie
CFLAGS_RELEASE += $(CFLAGS_LTO_PIE) -fsanitize=cfi,safe-stack
GO_LDFLAGS_RELEASE += "-w -s -linkmode=external -extldflags '$(LDFLAGS_RELEASE)'"

# Test with thread and memory sanitizers; needs associated libclang_rt libs.
test-race:
	CC=$(CC) CCLD=$(CCLD) CGO_CFLAGS="$(CFLAGS)" $(GO) test $(GO_BUILDFLAGS) -race -ldflags=$(GO_LDFLAGS) -coverpkg=.,./entropy .

# test-msan does not work on alpine (its compiler-rt lacks msan)
# but it works on fedora and void-musl.
test-msan:
	CC=clang CCLD=lld CGO_CFLAGS="$(CFLAGS_CFI)" $(GO) test $(GO_BUILDFLAGS) -buildmode=pie -msan -ldflags=$(GO_LDFLAGS_CFI) .

test-san: test-race test-msan

build-pie:
	CC=$(CC) CCLD=$(CCLD) CGO_CFLAGS="$(CFLAGS_CFI)" $(GO) build $(GO_BUILDFLAGS) -buildmode=pie -ldflags=$(GO_LDFLAGS_CFI) -o $(BIN) ./cmd/moac

# build-release builds a static-pie binary with sanitizers for CFI and either
# safe-stack (x86_64) or shadow-call-stack (ARMv8)
# the below should be run on a musl-based toolchain; works on Alpine or Void-musl
# Tends to cause crashes when linking with glibc
build-release:
	CC=clang CCLD=lld CGO_CFLAGS="$(CFLAGS_RELEASE)" $(GO) build $(GO_BUILDFLAGS) -buildmode=pie -ldflags=$(GO_LDFLAGS_RELEASE) -o $(BIN) ./cmd/moac

.PHONY: all lint test test-race test-msan test-san test-cov build build-release build-pie clean
